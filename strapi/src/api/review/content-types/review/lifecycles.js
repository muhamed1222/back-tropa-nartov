'use strict';

/**
 * Lifecycle hooks для Review (Отзыв)
 * - Автоматическая установка даты при создании
 */

module.exports = {
  /**
   * Перед созданием отзыва
   */
  async beforeCreate(event) {
    const { data } = event.params;
    
    // Автоматически устанавливаем дату, если не указана
    if (!data.date) {
      data.date = new Date();
    }
  },

  /**
   * Перед обновлением отзыва
   */
  async beforeUpdate(event) {
    // При обновлении дату не меняем
  },
};

