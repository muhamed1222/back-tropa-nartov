/**
 * Скрипт для добавления маршрутов по достопримечательностям Нальчика в Strapi
 * 
 * Использование:
 * node scripts/add-nalchik-routes.js
 * 
 * Требования:
 * - Достопримечательности Нальчика должны быть уже добавлены
 * - Настроены права доступа для Route (find, findOne, create)
 */

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

// Маршруты по Нальчику
const routes = [
  {
    name: 'Пешеходный маршрут "Центр Нальчика"',
    description: 'Увлекательная пешая прогулка по центру Нальчика. Маршрут включает посещение главных парков города, площадей и памятников. Идеально подходит для первого знакомства с городом.',
    route_type: 'Пеший', // Будет искаться по имени
    places: [
      'Курортный парк "Долина нарзанов"',
      'Площадь Абхазии',
      'Сквер имени Лермонтова',
      'Памятник "Вечная слава"',
      'Атажукинский сад',
    ],
    rating: 4.5,
  },
  {
    name: 'Культурный маршрут "Искусство и культура Нальчика"',
    description: 'Маршрут для любителей искусства и культуры. Посетите главные культурные достопримечательности Нальчика: театр, музеи, мечеть. Погрузитесь в культурную жизнь столицы Кабардино-Балкарии.',
    route_type: 'Авто',
    places: [
      'Кабардинский государственный драматический театр им. А. Шогенцукова',
      'Музей изобразительных искусств Кабардино-Балкарской Республики',
      'Соборная мечеть Нальчика',
    ],
    rating: 4.6,
  },
  {
    name: 'Исторический маршрут "Память поколений"',
    description: 'Маршрут по историческим и мемориальным местам Нальчика. Посетите памятники, посвященные истории города и героям Великой Отечественной войны.',
    route_type: 'Пеший',
    places: [
      'Памятник "Вечная слава"',
      'Памятник Ленину',
      'Сквер имени Лермонтова',
    ],
    rating: 4.7,
  },
  {
    name: 'Парковый маршрут "Зеленые легкие Нальчика"',
    description: 'Маршрут по всем паркам и скверам Нальчика. Идеально для любителей природы и спокойного отдыха. Прогуляйтесь по самым красивым зеленым зонам города.',
    route_type: 'Пеший',
    places: [
      'Курортный парк "Долина нарзанов"',
      'Атажукинский сад',
      'Сквер имени Лермонтова',
      'Городской парк культуры и отдыха',
    ],
    rating: 4.4,
  },
  {
    name: 'Комплексный маршрут "Знакомство с Нальчиком"',
    description: 'Полный маршрут для знакомства с Нальчиком за один день. Включает посещение главных достопримечательностей: парков, культурных объектов, памятников и площадей. Рекомендуется для туристов, впервые посещающих город.',
    route_type: 'Авто',
    places: [
      'Курортный парк "Долина нарзанов"',
      'Соборная мечеть Нальчика',
      'Кабардинский государственный драматический театр им. А. Шогенцукова',
      'Площадь Абхазии',
      'Памятник "Вечная слава"',
      'Музей изобразительных искусств Кабардино-Балкарской Республики',
      'Атажукинский сад',
    ],
    rating: 4.6,
  },
];

/**
 * Получить все места из API
 */
async function getPlaces() {
  try {
    const response = await axios.get(`${API_URL}/places?pagination[limit]=100`);
    return response.data.data || [];
  } catch (error) {
    console.log(`${colors.red}Ошибка получения мест: ${error.message}${colors.reset}`);
    return [];
  }
}

/**
 * Получить все типы маршрутов из API
 */
async function getRouteTypes() {
  try {
    const response = await axios.get(`${API_URL}/route-types`);
    return response.data.data || [];
  } catch (error) {
    console.log(`${colors.red}Ошибка получения типов маршрутов: ${error.message}${colors.reset}`);
    return [];
  }
}

/**
 * Найти место по имени
 */
function findPlaceByName(places, name) {
  return places.find(p => p.attributes.name === name);
}

/**
 * Найти тип маршрута по имени
 */
function findRouteTypeByName(routeTypes, name) {
  return routeTypes.find(rt => rt.attributes.name === name);
}

/**
 * Создание маршрута через Strapi API
 */
async function createRoute(routeData, places, routeTypes) {
  try {
    // Находим ID мест
    const placeIds = routeData.places
      .map(placeName => {
        const place = findPlaceByName(places, placeName);
        return place ? place.id : null;
      })
      .filter(id => id !== null);

    if (placeIds.length === 0) {
      return {
        success: false,
        error: 'Не найдены места для маршрута',
      };
    }

    // Находим тип маршрута
    const routeType = findRouteTypeByName(routeTypes, routeData.route_type);
    if (!routeType) {
      return {
        success: false,
        error: `Тип маршрута "${routeData.route_type}" не найден`,
      };
    }

    // Подготавливаем данные
    const data = {
      name: routeData.name,
      description: routeData.description,
      route_type: routeType.id,
      places: placeIds,
      is_active: true,
      rating: routeData.rating || 3.5,
    };

    const response = await axios.post(`${API_URL}/routes`, {
      data: data,
    }, {
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (response.data && response.data.data) {
      // Публикуем маршрут
      const routeId = response.data.data.id;
      await axios.put(
        `${API_URL}/routes/${routeId}`,
        {
          data: {
            ...data,
            publishedAt: new Date().toISOString(),
          },
        },
        {
          headers: {
            'Content-Type': 'application/json',
          },
        }
      );
      
      return { success: true, id: routeId, name: routeData.name };
    }
    
    return { success: false, error: 'Неожиданный ответ от API' };
  } catch (error) {
    return {
      success: false,
      error: error.response?.data?.error?.message || error.message,
      status: error.response?.status,
      details: error.response?.data,
    };
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
      const apiResponse = await axios.get(`${API_URL}/routes`, { timeout: 5000 });
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
  console.log(`${colors.blue}🚀 Запуск скрипта добавления маршрутов по Нальчику...${colors.reset}\n`);

  // Проверяем подключение к Strapi
  console.log(`${colors.yellow}Проверка подключения к Strapi...${colors.reset}`);
  const isConnected = await checkStrapiConnection();
  
  if (!isConnected) {
    console.log(`${colors.red}❌ Не удалось подключиться к Strapi по адресу ${STRAPI_URL}${colors.reset}`);
    console.log(`${colors.yellow}Убедитесь, что Strapi запущен: npm run develop${colors.reset}\n`);
    process.exit(1);
  }
  
  console.log(`${colors.green}✅ Strapi доступен${colors.reset}\n`);

  // Получаем места и типы маршрутов
  console.log(`${colors.yellow}Загрузка данных...${colors.reset}`);
  const [places, routeTypes] = await Promise.all([
    getPlaces(),
    getRouteTypes(),
  ]);

  if (places.length === 0) {
    console.log(`${colors.red}❌ Не найдены места. Сначала добавьте достопримечательности Нальчика.${colors.reset}`);
    console.log(`${colors.yellow}Запустите: node scripts/add-nalchik-places.js${colors.reset}\n`);
    process.exit(1);
  }

  if (routeTypes.length === 0) {
    console.log(`${colors.red}❌ Не найдены типы маршрутов. Добавьте типы маршрутов.${colors.reset}`);
    console.log(`${colors.yellow}Запустите: sqlite3 .tmp/data.db < scripts/add-route-types.sql${colors.reset}\n`);
    process.exit(1);
  }

  console.log(`${colors.green}✅ Загружено мест: ${places.length}, типов маршрутов: ${routeTypes.length}${colors.reset}\n`);

  // Добавляем маршруты
  console.log(`${colors.blue}Добавление маршрутов...${colors.reset}\n`);
  
  let successCount = 0;
  let failCount = 0;
  
  for (const route of routes) {
    console.log(`${colors.yellow}Добавляю: ${route.name}...${colors.reset}`);
    
    const result = await createRoute(route, places, routeTypes);
    
    if (result.success) {
      console.log(`${colors.green}✅ Добавлено: ${route.name} (ID: ${result.id})${colors.reset}`);
      console.log(`${colors.blue}   Мест в маршруте: ${route.places.length}${colors.reset}\n`);
      successCount++;
    } else {
      console.log(`${colors.red}❌ Ошибка: ${route.name}${colors.reset}`);
      if (result.status === 401 || result.status === 403) {
        console.log(`${colors.yellow}   Требуется авторизация. Настройте права доступа:${colors.reset}`);
        console.log(`${colors.yellow}   Settings → Users & Permissions → Roles → Public → Route → find/create${colors.reset}`);
      } else if (result.status === 405) {
        console.log(`${colors.yellow}   Метод не разрешен. Настройте права доступа:${colors.reset}`);
        console.log(`${colors.yellow}   Settings → Users & Permissions → Roles → Public → Route → create${colors.reset}`);
      } else {
        console.log(`${colors.red}   ${result.error} (статус: ${result.status || 'N/A'})${colors.reset}`);
        if (result.details) {
          console.log(`${colors.red}   Детали: ${JSON.stringify(result.details, null, 2)}${colors.reset}`);
        }
        console.log('');
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
    console.log(`${colors.green}🎉 Маршруты по Нальчику успешно добавлены!${colors.reset}`);
    console.log(`${colors.blue}Проверьте в админ-панели: ${STRAPI_URL}/admin${colors.reset}\n`);
  }
}

// Запуск
main().catch(error => {
  console.error(`${colors.red}Критическая ошибка:${colors.reset}`, error);
  process.exit(1);
});

