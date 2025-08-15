-- Create patients table following FHIR Patient resource structure
CREATE TABLE IF NOT EXISTS patients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    identifier JSONB NOT NULL DEFAULT '[]'::jsonb,
    active BOOLEAN NOT NULL DEFAULT true,
    name JSONB NOT NULL DEFAULT '[]'::jsonb,
    telecom JSONB DEFAULT '[]'::jsonb,
    gender VARCHAR(20),
    birth_date DATE,
    deceased_boolean BOOLEAN DEFAULT false,
    deceased_date_time TIMESTAMP WITH TIME ZONE,
    address JSONB DEFAULT '[]'::jsonb,
    marital_status JSONB,
    multiple_birth_boolean BOOLEAN,
    multiple_birth_integer INTEGER,
    photo JSONB DEFAULT '[]'::jsonb,
    contact JSONB DEFAULT '[]'::jsonb,
    communication JSONB DEFAULT '[]'::jsonb,
    general_practitioner JSONB DEFAULT '[]'::jsonb,
    managing_organization JSONB,
    link JSONB DEFAULT '[]'::jsonb,
    meta JSONB DEFAULT '{}'::jsonb,
    implicit_rules TEXT,
    language VARCHAR(10),
    text JSONB,
    contained JSONB DEFAULT '[]'::jsonb,
    extension JSONB DEFAULT '[]'::jsonb,
    modifier_extension JSONB DEFAULT '[]'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    version INTEGER DEFAULT 1
);

-- Create indexes for performance
CREATE INDEX idx_patients_identifier ON patients USING GIN (identifier);
CREATE INDEX idx_patients_name ON patients USING GIN (name);
CREATE INDEX idx_patients_active ON patients (active);
CREATE INDEX idx_patients_gender ON patients (gender);
CREATE INDEX idx_patients_birth_date ON patients (birth_date);
CREATE INDEX idx_patients_created_at ON patients (created_at);
CREATE INDEX idx_patients_updated_at ON patients (updated_at);

-- Create trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    NEW.version = OLD.version + 1;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_patients_updated_at 
    BEFORE UPDATE ON patients 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
