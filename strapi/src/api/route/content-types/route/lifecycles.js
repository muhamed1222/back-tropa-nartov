'use strict';

/**
 * Lifecycle hooks для автоматического расчета данных маршрута
 * - distance_km: расстояние между местами (формула Хаверсина)
 * - duration_hours: время прохождения (расчетное, зависит от типа маршрута)
 * - places_count: количество мест в маршруте
 * 
 * Использует компонент route.route-stop для сохранения порядка мест
 */

// Радиус Земли в километрах
const R = 6371;

// Скорости по умолчанию для разных типов маршрутов (км/ч)
const DEFAULT_SPEEDS = {
  walking: 3,      // Пеший туризм
  cycling: 12,     // Велосипед
  driving: 50,     // Автомобиль
  default: 3,      // По умолчанию
};

/**
 * Преобразование градусов в радианы
 */
function toRad(degrees) {
  return degrees * (Math.PI / 180);
}

/**
 * Расчет расстояния между двумя точками по формуле Хаверсина
 * @param {Object} p1 - первая точка {latitude, longitude}
 * @param {Object} p2 - вторая точка {latitude, longitude}
 * @returns {number} расстояние в километрах
 */
function getDistance(p1, p2) {
  const dLat = toRad(p2.latitude - p1.latitude);
  const dLng = toRad(p2.longitude - p1.longitude);
  const lat1 = toRad(p1.latitude);
  const lat2 = toRad(p2.latitude);

  const a =
    Math.sin(dLat / 2) ** 2 +
    Math.sin(dLng / 2) ** 2 * Math.cos(lat1) * Math.cos(lat2);

  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  return R * c;
}

/**
 * Получить среднюю скорость для типа маршрута
 * @param {number} typeId - ID типа маршрута
 * @returns {Promise<number>} средняя скорость в км/ч
 */
async function getAverageSpeed(typeId) {
  if (!typeId) {
    return DEFAULT_SPEEDS.default;
  }

  try {
    const routeType = await strapi.entityService.findOne('api::type.type', typeId, {
      fields: ['id', 'name'],
    });

    if (!routeType) {
      return DEFAULT_SPEEDS.default;
    }

    const typeName = routeType.name?.toLowerCase() || '';
    
    // Определяем скорость по названию типа
    if (typeName.includes('пеший') || typeName.includes('ходьба') || typeName.includes('walking')) {
      return DEFAULT_SPEEDS.walking;
    }
    if (typeName.includes('вело') || typeName.includes('велосипед') || typeName.includes('cycling')) {
      return DEFAULT_SPEEDS.cycling;
    }
    if (typeName.includes('авто') || typeName.includes('машина') || typeName.includes('driving')) {
      return DEFAULT_SPEEDS.driving;
    }

    return DEFAULT_SPEEDS.default;
  } catch (error) {
    console.warn(`Не удалось получить тип маршрута ${typeId}:`, error.message);
    return DEFAULT_SPEEDS.default;
  }
}

/**
 * Валидация остановок маршрута
 * @param {Array} stops - массив остановок (компоненты route-stop)
 * @throws {Error} если валидация не пройдена
 */
function validateStops(stops) {
  if (!stops || !Array.isArray(stops)) {
    return; // Остановки не обязательны
  }

  if (stops.length === 0) {
    return; // Пустой маршрут допустим
  }

  // Проверка минимального количества мест (если есть места, должно быть минимум 2)
  if (stops.length === 1) {
    throw new Error('Маршрут должен содержать минимум 2 места для расчета расстояния');
  }

  // Проверка на дубликаты мест
  const placeIds = stops
    .map((stop) => {
      // stop может быть объектом с place.id или просто ID
      if (typeof stop === 'object' && stop.place) {
        return typeof stop.place === 'object' ? stop.place.id : stop.place;
      }
      return null;
    })
    .filter((id) => id != null);

  const uniquePlaceIds = new Set(placeIds);
  if (placeIds.length !== uniquePlaceIds.size) {
    throw new Error('В маршруте не должно быть дублирующихся мест');
  }

  // Автоматическая установка order, если не указан
  stops.forEach((stop, index) => {
    if (typeof stop === 'object' && (stop.order === undefined || stop.order === null)) {
      stop.order = index;
    }
  });

  // Сортировка по order
  stops.sort((a, b) => {
    const aOrder = typeof a === 'object' ? (a.order ?? 0) : 0;
    const bOrder = typeof b === 'object' ? (b.order ?? 0) : 0;
    return aOrder - bOrder;
  });
}

/**
 * Расчет данных маршрута (расстояние, время, количество мест)
 * @param {Object} data - данные маршрута
 * @param {Array} data.stops - массив остановок (компоненты route-stop)
 * @param {number} data.type_id - ID типа маршрута
 * @returns {Promise<Object>} обновленные данные с рассчитанными полями
 */
async function calculateRouteData(data) {
  // Валидация остановок
  try {
    validateStops(data.stops);
  } catch (error) {
    throw new Error(`Ошибка валидации маршрута: ${error.message}`);
  }

  // Если остановки не указаны, очищаем рассчитанные поля
  if (!data.stops || !Array.isArray(data.stops) || data.stops.length === 0) {
    data.distance_km = 0;
    data.duration_hours = 0;
    data.places_count = 0;
    return data;
  }

  try {
    // Извлекаем ID мест из компонентов stops
    const placeIds = data.stops
      .map((stop) => {
        if (typeof stop === 'object' && stop.place) {
          return typeof stop.place === 'object' ? stop.place.id : stop.place;
        }
        return null;
      })
      .filter((id) => id != null);

    if (placeIds.length === 0) {
      data.distance_km = 0;
      data.duration_hours = 0;
      data.places_count = 0;
      return data;
    }

    // Загружаем места с координатами
    const places = await strapi.entityService.findMany('api::place.place', {
      filters: {
        id: {
          $in: placeIds,
        },
      },
      fields: ['id', 'latitude', 'longitude', 'name'],
    });

    // Проверяем, что все места найдены
    if (places.length !== placeIds.length) {
      const foundIds = places.map((p) => p.id);
      const missingIds = placeIds.filter((id) => !foundIds.includes(id));
      throw new Error(
        `Некоторые места не найдены: ${missingIds.join(', ')}. Проверьте существование этих мест.`
      );
    }

    // Сортируем места по порядку из stops (по order)
    const sortedStops = [...data.stops].sort((a, b) => {
      const aOrder = typeof a === 'object' ? (a.order ?? 0) : 0;
      const bOrder = typeof b === 'object' ? (b.order ?? 0) : 0;
      return aOrder - bOrder;
    });

    const sortedPlaces = sortedStops
      .map((stop) => {
        const placeId = typeof stop === 'object' && stop.place
          ? (typeof stop.place === 'object' ? stop.place.id : stop.place)
          : null;
        return places.find((p) => p.id === placeId);
      })
      .filter((p) => p != null);

    // Проверяем, что у всех мест есть координаты
    const placesWithoutCoords = sortedPlaces.filter(
      (p) => p.latitude == null || p.longitude == null || isNaN(p.latitude) || isNaN(p.longitude)
    );

    if (placesWithoutCoords.length > 0) {
      const placeNames = placesWithoutCoords.map((p) => p.name || `ID ${p.id}`).join(', ');
      throw new Error(
        `У следующих мест отсутствуют координаты: ${placeNames}. Пожалуйста, укажите координаты (latitude, longitude) для всех мест маршрута.`
      );
    }

    if (sortedPlaces.length === 0) {
      data.distance_km = 0;
      data.duration_hours = 0;
      data.places_count = 0;
      return data;
    }

    // Рассчитываем общее расстояние между последовательными точками
    let totalDistance = 0;
    for (let i = 0; i < sortedPlaces.length - 1; i++) {
      const p1 = sortedPlaces[i];
      const p2 = sortedPlaces[i + 1];
      const dist = getDistance(
        { latitude: parseFloat(p1.latitude), longitude: parseFloat(p1.longitude) },
        { latitude: parseFloat(p2.latitude), longitude: parseFloat(p2.longitude) }
      );
      totalDistance += dist;
    }

    // Сохраняем рассчитанные значения
    data.distance_km = parseFloat(totalDistance.toFixed(2));
    data.places_count = sortedPlaces.length;

    // Расчет времени прохождения (динамическая скорость)
    const averageSpeed = await getAverageSpeed(data.type_id);
    const duration = totalDistance / averageSpeed;
    data.duration_hours = parseFloat(duration.toFixed(1));

    return data;
  } catch (error) {
    console.error('Ошибка при расчете данных маршрута:', error);
    // Пробрасываем ошибку валидации дальше
    if (error.message && error.message.includes('Ошибка валидации') || error.message.includes('отсутствуют координаты') || error.message.includes('не найдены')) {
      throw error;
    }
    // В случае другой ошибки сохраняем значения по умолчанию
    data.distance_km = 0;
    data.duration_hours = 0;
    data.places_count = data.stops?.length || 0;
    return data;
  }
}

module.exports = {
  /**
   * Перед созданием маршрута
   */
  async beforeCreate(event) {
    const { data } = event.params;
    await calculateRouteData(data);
  },

  /**
   * Перед обновлением маршрута
   */
  async beforeUpdate(event) {
    const { data } = event.params;

    // Всегда пересчитываем, если обновляются остановки
    if (data.stops !== undefined) {
      await calculateRouteData(data);
    } else {
      // Если остановки не меняются, но нужно пересчитать данные для существующего маршрута
      // (например, если изменились координаты места)
      try {
        const routeId = event.params.where?.id || event.params.where;
        if (routeId) {
          const existingRoute = await strapi.entityService.findOne(
            'api::route.route',
            routeId,
            {
              populate: ['stops.place'],
            }
          );

          if (existingRoute && existingRoute.stops && existingRoute.stops.length > 0) {
            // Используем существующие остановки для пересчета
            const tempData = {
              ...data,
              stops: existingRoute.stops,
              type_id: data.type_id || existingRoute.type_id,
            };
            await calculateRouteData(tempData);

            // Объединяем рассчитанные значения с обновляемыми данными
            Object.assign(data, {
              distance_km: tempData.distance_km,
              duration_hours: tempData.duration_hours,
              places_count: tempData.places_count,
            });
          }
        }
      } catch (error) {
        console.error('Ошибка при загрузке существующего маршрута для пересчета:', error);
        // Если не удалось загрузить, просто пропускаем пересчет
      }
    }
  },

  /**
   * После создания/обновления маршрута
   */
  async afterCreate(event) {
    // Можно добавить дополнительную логику, например, кеширование
  },

  async afterUpdate(event) {
    // Можно добавить дополнительную логику, например, кеширование
  },
};
