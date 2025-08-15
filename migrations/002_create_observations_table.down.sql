-- Drop observations table and related objects
DROP TRIGGER IF EXISTS update_observations_updated_at ON observations;
DROP TABLE IF EXISTS observations;
