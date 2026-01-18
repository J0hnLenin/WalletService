SET session_replication_role = 'replica';

DO $$
DECLARE
    bucket_id INT;
BEGIN
    FOR bucket_id IN 0..63 LOOP
        EXECUTE format('DROP SCHEMA IF EXISTS bucket_%s CASCADE', bucket_id);
    END LOOP;
END $$;

SET session_replication_role = 'origin';

DO $$
DECLARE
    bucket_id INT;
    schema_name TEXT;
BEGIN
    FOR bucket_id IN 0..31 LOOP
        schema_name := 'bucket_' || bucket_id;
        
        EXECUTE format('CREATE SCHEMA IF NOT EXISTS %I', schema_name);
        
        EXECUTE format('
            CREATE TABLE IF NOT EXISTS %I.wallets (
                id UUID PRIMARY KEY,
                amount INT NOT NULL CHECK (amount > 0),
                created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
            )', schema_name);
    END LOOP;
END $$;
