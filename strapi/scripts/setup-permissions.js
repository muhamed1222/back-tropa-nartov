/**
 * Скрипт для настройки прав доступа в Strapi
 * Настраивает права для Public роли, чтобы можно было создавать места через API
 * 
 * Использование:
 * node scripts/setup-permissions.js
 */

const axios = require('axios');

const STRAPI_URL = 'http://localhost:1337';
const API_URL = `${STRAPI_URL}/api`;

// Цвета для вывода в консоль
const colors = {
  reset: '\x1b[0m',
  green: '\x1b[32m',
  red: '\x1b[31m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
};

/**
 * Получить токен администратора
 * Примечание: В продакшене используйте API Token из Settings
 */
async function getAdminToken() {
  // Для разработки можно использовать API Token
  // Или попросить пользователя войти в админку и скопировать токен
  console.log(`${colors.yellow}⚠️  Для настройки прав нужен административный доступ${colors.reset}`);
  console.log(`${colors.blue}Вариант 1: Настройте права вручную через админку:${colors.reset}`);
  console.log(`${colors.blue}  Settings → Users & Permissions → Roles → Public → Place${colors.reset}`);
  console.log(`${colors.blue}  Включите: find, findOne, create${colors.reset}\n`);
  
  // Альтернатива: можно использовать API Token
  const apiToken = process.env.STRAPI_API_TOKEN;
  if (apiToken) {
    return apiToken;
  }
  
  return null;
}

/**
 * Настройка прав доступа через API (если есть токен)
 */
async function setupPermissions(token) {
  if (!token) {
    return false;
  }

  try {
    // Получаем роль Public
    const rolesResponse = await axios.get(
      `${API_URL}/users-permissions/roles`,
      {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      }
    );

    const publicRole = rolesResponse.data.roles.find(role => role.type === 'public');
    
    if (!publicRole) {
      console.log(`${colors.red}❌ Роль Public не найдена${colors.reset}`);
      return false;
    }

    // Обновляем права для Place
    const updatedPermissions = {
      ...publicRole.permissions,
      'api::place.place': {
        controllers: {
          place: {
            find: { enabled: true },
            findOne: { enabled: true },
            create: { enabled: true },
          },
        },
      },
    };

    await axios.put(
      `${API_URL}/users-permissions/roles/${publicRole.id}`,
      {
        permissions: updatedPermissions,
      },
      {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      }
    );

    console.log(`${colors.green}✅ Права доступа настроены${colors.reset}`);
    return true;
  } catch (error) {
    console.log(`${colors.red}❌ Ошибка настройки прав: ${error.message}${colors.reset}`);
    return false;
  }
}

/**
 * Проверка доступности Strapi
 */
async function checkStrapiConnection() {
  try {
    const response = await axios.get(`${STRAPI_URL}/admin`, { timeout: 5000 });
    return response.status === 200;
  } catch (error) {
    try {
      const apiResponse = await axios.get(`${API_URL}/places`, { timeout: 5000 });
      return true;
    } catch (apiError) {
      return false;
    }
  }
}

/**
 * Основная функция
 */
async function main() {
  console.log(`${colors.blue}🚀 Настройка прав доступа для добавления мест...${colors.reset}\n`);

  // Проверяем подключение к Strapi
  console.log(`${colors.yellow}Проверка подключения к Strapi...${colors.reset}`);
  const isConnected = await checkStrapiConnection();
  
  if (!isConnected) {
    console.log(`${colors.red}❌ Не удалось подключиться к Strapi по адресу ${STRAPI_URL}${colors.reset}`);
    console.log(`${colors.yellow}Убедитесь, что Strapi запущен: npm run develop${colors.reset}\n`);
    process.exit(1);
  }
  
  console.log(`${colors.green}✅ Strapi доступен${colors.reset}\n`);

  // Пробуем настроить права
  const token = await getAdminToken();
  const permissionsSet = await setupPermissions(token);

  if (!permissionsSet) {
    console.log(`\n${colors.yellow}📝 Настройте права вручную:${colors.reset}`);
    console.log(`1. Откройте: ${STRAPI_URL}/admin`);
    console.log(`2. Перейдите: Settings → Users & Permissions → Roles → Public`);
    console.log(`3. Найдите "Place" в списке`);
    console.log(`4. Включите права: find, findOne, create`);
    console.log(`5. Нажмите "Save"`);
    console.log(`6. Затем запустите: node scripts/add-places.js\n`);
    process.exit(0);
  }

  console.log(`\n${colors.green}🎉 Права доступа настроены!${colors.reset}`);
  console.log(`${colors.blue}Теперь можно запустить: node scripts/add-places.js${colors.reset}\n`);
}

// Запуск
main().catch(error => {
  console.error(`${colors.red}Критическая ошибка:${colors.reset}`, error);
  process.exit(1);
});

