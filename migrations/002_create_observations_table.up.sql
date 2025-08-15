-- Create observations table following FHIR Observation resource structure
CREATE TABLE IF NOT EXISTS observations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    identifier JSONB DEFAULT '[]'::jsonb,
    based_on JSONB DEFAULT '[]'::jsonb,
    part_of JSONB DEFAULT '[]'::jsonb,
    status VARCHAR(50) NOT NULL CHECK (status IN ('registered', 'preliminary', 'final', 'amended', 'corrected', 'cancelled', 'entered-in-error', 'unknown')),
    category JSONB DEFAULT '[]'::jsonb,
    code JSONB NOT NULL,
    subject JSONB NOT NULL,
    focus JSONB DEFAULT '[]'::jsonb,
    encounter JSONB,
    effective_date_time TIMESTAMP WITH TIME ZONE,
    effective_period JSONB,
    effective_timing JSONB,
    effective_instant TIMESTAMP WITH TIME ZONE,
    issued TIMESTAMP WITH TIME ZONE,
    performer JSONB DEFAULT '[]'::jsonb,
    value_quantity JSONB,
    value_codeable_concept JSONB,
    value_string TEXT,
    value_boolean BOOLEAN,
    value_integer INTEGER,
    value_range JSONB,
    value_ratio JSONB,
    value_sampled_data JSONB,
    value_time TIME,
    value_date_time TIMESTAMP WITH TIME ZONE,
    value_period JSONB,
    data_absent_reason JSONB,
    interpretation JSONB DEFAULT '[]'::jsonb,
    note JSONB DEFAULT '[]'::jsonb,
    body_site JSONB,
    method JSONB,
    specimen JSONB,
    device JSONB,
    reference_range JSONB DEFAULT '[]'::jsonb,
    has_member JSONB DEFAULT '[]'::jsonb,
    derived_from JSONB DEFAULT '[]'::jsonb,
    component JSONB DEFAULT '[]'::jsonb,
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
CREATE INDEX idx_observations_identifier ON observations USING GIN (identifier);
CREATE INDEX idx_observations_status ON observations (status);
CREATE INDEX idx_observations_category ON observations USING GIN (category);
CREATE INDEX idx_observations_code ON observations USING GIN (code);
CREATE INDEX idx_observations_subject ON observations USING GIN (subject);
CREATE INDEX idx_observations_effective_date_time ON observations (effective_date_time);
CREATE INDEX idx_observations_issued ON observations (issued);
CREATE INDEX idx_observations_created_at ON observations (created_at);
CREATE INDEX idx_observations_updated_at ON observations (updated_at);

-- Create trigger for updated_at
CREATE TRIGGER update_observations_updated_at 
    BEFORE UPDATE ON observations 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
