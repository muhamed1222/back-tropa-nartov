// config/admin.js
module.exports = ({ env }) => ({
  auth: {
    secret: env('ADMIN_JWT_SECRET', 'someSecretKey'),
  },
  apiToken: {
    salt: env('API_TOKEN_SALT', 'someSalt'),
  },
  transfer: {
    token: {
      salt: env('TRANSFER_TOKEN_SALT', 'someTransferSalt'),
    },
  },
  flags: {
    nps: env.bool('FLAG_NPS', true),
    promoteEE: env.bool('FLAG_PROMOTE_EE', true),
  },
  // ВАЖНО: Эта конфигурация должна быть на верхнем уровне
  i18n: {
    locales: ['en', 'ru'],
    defaultLocale: 'en'
  }
});