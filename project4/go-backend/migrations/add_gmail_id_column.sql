-- Add gmail_id column to emails table if it doesn't exist
DO $$
BEGIN
    -- Check if the column already exists
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name='emails' 
        AND column_name='gmail_id'
    ) THEN
        -- Add the gmail_id column
        ALTER TABLE emails ADD COLUMN gmail_id VARCHAR(255);
        
        -- Add unique constraint
        ALTER TABLE emails ADD CONSTRAINT emails_gmail_id_unique UNIQUE (gmail_id);
        
        -- Add index for better performance
        CREATE INDEX IF NOT EXISTS idx_emails_gmail_id ON emails(gmail_id);
        
        RAISE NOTICE 'gmail_id column added to emails table';
    ELSE
        RAISE NOTICE 'gmail_id column already exists in emails table';
    END IF;
END $$;
