-- Drop patients table and related objects
DROP TRIGGER IF EXISTS update_patients_updated_at ON patients;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS patients;
