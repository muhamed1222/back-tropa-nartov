/**
 * Скрипт для добавления тестовых мест в Strapi
 * 
 * Использование:
 * node scripts/add-places.js
 */

// Используем встроенный модуль fetch или axios
let axios;
try {
  axios = require('axios');
} catch (error) {
  console.log('Устанавливаем axios...');
  console.log('Запустите: npm install axios');
  process.exit(1);
}

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

// Тестовые места для добавления
const places = [
  {
    name: 'Парк Атажукинский',
    short_description: 'Крупнейший парк Нальчика площадью 250 гектаров. Один из старейших парков Кабардино-Балкарии, заложен в 1847 году.',
    address: 'г. Нальчик, ул. Толстого, 2',
    working_hours: '06:00-24:00, без выходных',
    contacts_phone: '+78662225566',
    contacts_website: 'https://park-nalchik.ru',
    latitude: 43.4925,
    longitude: 43.6123,
    history: 'Парк был заложен в 1847 году по инициативе генерала Атажукина. За годы существования парк неоднократно реконструировался, сохраняя свою уникальную природную красоту.',
    is_active: true,
  },
  {
    name: 'Национальный музей КБР',
    short_description: 'Главный музей Кабардино-Балкарской Республики, рассказывающий об истории и культуре народов региона.',
    address: 'г. Нальчик, ул. Горького, 62',
    working_hours: '09:00-18:00, выходной: понедельник',
    contacts_phone: '+78662773544',
    contacts_website: 'https://museum-kbr.ru',
    latitude: 43.4833,
    longitude: 43.6017,
    history: 'Музей основан в 1921 году. Содержит уникальные коллекции археологических находок, этнографических экспонатов и произведений искусства.',
    is_active: true,
  },
  {
    name: 'Гора Эльбрус',
    short_description: 'Высочайшая вершина России и Европы. Высота 5642 метра. Популярное место для альпинизма и туризма.',
    address: 'Кабардино-Балкария, Приэльбрусье',
    working_hours: 'Круглосуточно',
    contacts_phone: '+78663872222',
    contacts_website: 'https://elbrus.ru',
    latitude: 43.3550,
    longitude: 42.4392,
    history: 'Эльбрус - это потухший вулкан, высочайшая вершина Кавказа. Первое восхождение было совершено в 1829 году. Сейчас это популярный горнолыжный курорт и место для альпинизма.',
    is_active: true,
  },
  {
    name: 'Чегемские водопады',
    short_description: 'Красивые водопады в ущелье реки Чегем. Состоят из нескольких каскадов, самый высокий - Су-Аузу.',
    address: 'Кабардино-Балкария, Чегемское ущелье',
    working_hours: '08:00-20:00, без выходных',
    contacts_phone: '+78663876666',
    latitude: 43.3167,
    longitude: 43.1833,
    history: 'Чегемские водопады - одна из главных природных достопримечательностей Кабардино-Балкарии. Водопады особенно красивы зимой, когда замерзают и образуют ледяные скульптуры.',
    is_active: true,
  },
  {
    name: 'Голубое озеро (Церик-Кель)',
    short_description: 'Уникальное карстовое озеро ярко-голубого цвета. Одно из самых глубоких озёр в России.',
    address: 'Кабардино-Балкария, Черекский район, село Бабугент',
    working_hours: '08:00-19:00',
    contacts_phone: '+78663875555',
    latitude: 43.2344,
    longitude: 43.5481,
    history: 'Голубое озеро - карстовое озеро, одно из самых глубоких в России (глубина более 250 метров). Вода имеет ярко-голубой цвет из-за содержания сероводорода. Озеро никогда не замерзает.',
    is_active: true,
  },
  {
    name: 'Баксанское ущелье',
    short_description: 'Живописное ущелье в Приэльбрусье, популярное место для туризма и экскурсий.',
    address: 'Кабардино-Балкария, Баксанский район',
    working_hours: 'Круглосуточно',
    latitude: 43.3833,
    longitude: 42.7500,
    history: 'Баксанское ущелье - одно из самых красивых ущелий Кавказа. Простирается от города Баксан до подножия Эльбруса. Известно своими живописными видами и развитой туристической инфраструктурой.',
    is_active: true,
  },
];

/**
 * Создание места через Strapi API
 */
async function createPlace(placeData) {
  try {
    const response = await axios.post(`${API_URL}/places`, {
      data: placeData,
    }, {
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (response.data && response.data.data) {
      // Публикуем место
      const placeId = response.data.data.id;
      await axios.put(
        `${API_URL}/places/${placeId}`,
        {
          data: {
            ...placeData,
            publishedAt: new Date().toISOString(),
          },
        },
        {
          headers: {
            'Content-Type': 'application/json',
          },
        }
      );
      
      return { success: true, id: placeId, name: placeData.name };
    }
    
    return { success: false, error: 'Неожиданный ответ от API' };
  } catch (error) {
    return {
      success: false,
      error: error.response?.data?.error?.message || error.message,
      status: error.response?.status,
    };
  }
}

/**
 * Проверка доступности Strapi
 */
async function checkStrapiConnection() {
  try {
    // Проверяем доступность через админку или API
    const response = await axios.get(`${STRAPI_URL}/admin`, { timeout: 5000 });
    return response.status === 200;
  } catch (error) {
    // Если админка не доступна, пробуем API
    try {
      const apiResponse = await axios.get(`${API_URL}/places`, { timeout: 5000 });
      return true; // Даже если 404, значит Strapi работает
    } catch (apiError) {
      return false;
    }
  }
}

/**
 * Основная функция
 */
async function main() {
  console.log(`${colors.blue}🚀 Запуск скрипта добавления мест...${colors.reset}\n`);

  // Проверяем подключение к Strapi
  console.log(`${colors.yellow}Проверка подключения к Strapi...${colors.reset}`);
  const isConnected = await checkStrapiConnection();
  
  if (!isConnected) {
    console.log(`${colors.red}❌ Не удалось подключиться к Strapi по адресу ${STRAPI_URL}${colors.reset}`);
    console.log(`${colors.yellow}Убедитесь, что Strapi запущен: npm run develop${colors.reset}\n`);
    process.exit(1);
  }
  
  console.log(`${colors.green}✅ Strapi доступен${colors.reset}\n`);

  // Добавляем места
  console.log(`${colors.blue}Добавление мест...${colors.reset}\n`);
  
  let successCount = 0;
  let failCount = 0;
  
  for (const place of places) {
    console.log(`${colors.yellow}Добавляю: ${place.name}...${colors.reset}`);
    
    const result = await createPlace(place);
    
    if (result.success) {
      console.log(`${colors.green}✅ Добавлено: ${place.name} (ID: ${result.id})${colors.reset}\n`);
      successCount++;
    } else {
      console.log(`${colors.red}❌ Ошибка: ${place.name}${colors.reset}`);
      if (result.status === 401 || result.status === 403) {
        console.log(`${colors.yellow}   Требуется авторизация. Настройте права доступа:${colors.reset}`);
        console.log(`${colors.yellow}   Запустите: node scripts/setup-permissions.js${colors.reset}`);
        console.log(`${colors.yellow}   Или настройте вручную в админке:${colors.reset}`);
        console.log(`${colors.yellow}   Settings → Users & Permissions → Roles → Public → Place → find/create${colors.reset}`);
      } else if (result.status === 405) {
        console.log(`${colors.yellow}   Метод не разрешен. Настройте права доступа:${colors.reset}`);
        console.log(`${colors.yellow}   Settings → Users & Permissions → Roles → Public → Place → create${colors.reset}`);
        console.log(`${colors.yellow}   Затем перезапустите Strapi${colors.reset}`);
      } else {
        console.log(`${colors.red}   ${result.error} (статус: ${result.status || 'N/A'})${colors.reset}\n`);
      }
      failCount++;
    }
    
    // Небольшая задержка между запросами
    await new Promise(resolve => setTimeout(resolve, 500));
  }

  // Итоги
  console.log(`\n${colors.blue}═══════════════════════════════════════${colors.reset}`);
  console.log(`${colors.green}✅ Успешно добавлено: ${successCount}${colors.reset}`);
  if (failCount > 0) {
    console.log(`${colors.red}❌ Ошибок: ${failCount}${colors.reset}`);
  }
  console.log(`${colors.blue}═══════════════════════════════════════${colors.reset}\n`);

  if (successCount > 0) {
    console.log(`${colors.green}🎉 Места успешно добавлены!${colors.reset}`);
    console.log(`${colors.blue}Проверьте в админ-панели: ${STRAPI_URL}/admin${colors.reset}\n`);
  }
}

// Запуск
main().catch(error => {
  console.error(`${colors.red}Критическая ошибка:${colors.reset}`, error);
  process.exit(1);
});

