module.exports = ({ env }) => {
  // По умолчанию используем SQLite для разработки
  const client = env('DATABASE_CLIENT') || 'sqlite';
  
  // Всегда используем SQLite, если не указан DATABASE_CLIENT=postgres
  if (client !== 'postgres') {
    // Используем SQLite для разработки (проще настроить)
    return {
      connection: {
        client: 'sqlite',
        connection: {
          filename: env('DATABASE_FILENAME', '.tmp/data.db'),
        },
        useNullAsDefault: true,
      },
    };
  }
  
  // Используем PostgreSQL для production
  return {
    connection: {
      client: 'postgres',
      connection: {
        host: env('DATABASE_HOST', 'localhost'),
        port: env.int('DATABASE_PORT', 5432),
        database: env('DATABASE_NAME', 'tropa_nartov'),
        user: env('DATABASE_USERNAME', 'postgres'),
        password: env('DATABASE_PASSWORD', 'password'),
        ssl: env.bool('DATABASE_SSL', false),
      },
      pool: {
        min: 0,
        max: 10,
      },
    },
  };
};

