SELECT 'CREATE DATABASE test1'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'test1')\gexec

SELECT 'CREATE DATABASE test2'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'test2')\gexec

SELECT 'CREATE DATABASE sonar'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'sonar')\gexec

DO $$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_roles WHERE rolname = 'sonar'
   ) THEN
      CREATE USER sonar WITH PASSWORD 'sonar';
   END IF;

   IF EXISTS (
      SELECT FROM pg_database WHERE datname = 'sonar'
   ) THEN
      ALTER DATABASE sonar OWNER TO sonar;
   END IF;

   IF EXISTS (
      SELECT FROM pg_database WHERE datname = 'sonar'
   ) AND EXISTS (
      SELECT FROM pg_roles WHERE rolname = 'sonar'
   ) THEN
      GRANT ALL PRIVILEGES ON DATABASE sonar TO sonar;
   END IF;
END
$$;