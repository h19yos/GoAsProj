ALTER TABLE module_info
    ADD CONSTRAINT check_updated_at_after_created_at
        CHECK (updated_at >= created_at);

ALTER TABLE module_info
    ADD CONSTRAINT check_module_duration_range
        CHECK (module_duration > 5 AND module_duration <= 15);
