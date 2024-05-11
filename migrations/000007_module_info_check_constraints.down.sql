ALTER TABLE module_info
    DROP CONSTRAINT IF EXISTS check_updated_at_after_created_at,
    DROP CONSTRAINT IF EXISTS check_module_duration_range;
