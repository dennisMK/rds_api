package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"healthcare-api/internal/database"
	"healthcare-api/internal/models"

	"github.com/google/uuid"
)

type PatientRepository struct {
	*BaseRepository
}

func NewPatientRepository(db *database.DB) *PatientRepository {
	return &PatientRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *PatientRepository) Create(ctx context.Context, patient *models.Patient) error {
	query := `
		INSERT INTO patients (
			id, identifier, active, name, telecom, gender, birth_date,
			deceased_boolean, deceased_date_time, address, marital_status,
			multiple_birth_boolean, multiple_birth_integer, photo, contact,
			communication, general_practitioner, managing_organization, link,
			meta, implicit_rules, language, text, contained, extension, modifier_extension
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26
		) RETURNING created_at, updated_at, version
	`

	err := r.db.QueryRowContext(ctx, query,
		patient.ID,
		toJSON(patient.Identifier),
		patient.Active,
		toJSON(patient.Name),
		toJSON(patient.Telecom),
		patient.Gender,
		patient.BirthDate,
		patient.DeceasedBoolean,
		patient.DeceasedDateTime,
		toJSON(patient.Address),
		toJSON(patient.MaritalStatus),
		patient.MultipleBirthBoolean,
		patient.MultipleBirthInteger,
		toJSON(patient.Photo),
		toJSON(patient.Contact),
		toJSON(patient.Communication),
		toJSON(patient.GeneralPractitioner),
		toJSON(patient.ManagingOrganization),
		toJSON(patient.Link),
		toJSON(patient.Meta),
		patient.ImplicitRules,
		patient.Language,
		toJSON(patient.Text),
		toJSON(patient.Contained),
		toJSON(patient.Extension),
		toJSON(patient.ModifierExtension),
	).Scan(&patient.CreatedAt, &patient.UpdatedAt, &patient.Version)

	if err != nil {
		return fmt.Errorf("failed to create patient: %w", err)
	}

	// Log audit trail
	auditLog := &AuditLog{
		ResourceType: "Patient",
		ResourceID:   patient.ID,
		Action:       "CREATE",
		NewValues:    mustMarshalJSON(patient),
	}
	
	if err := r.LogAudit(ctx, auditLog); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to log audit: %v\n", err)
	}

	return nil
}

func (r *PatientRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Patient, error) {
	query := `
		SELECT id, identifier, active, name, telecom, gender, birth_date,
			   deceased_boolean, deceased_date_time, address, marital_status,
			   multiple_birth_boolean, multiple_birth_integer, photo, contact,
			   communication, general_practitioner, managing_organization, link,
			   meta, implicit_rules, language, text, contained, extension, 
			   modifier_extension, created_at, updated_at, version
		FROM patients WHERE id = $1
	`

	patient := &models.Patient{}
	var identifier, name, telecom, address, maritalStatus, photo, contact []byte
	var communication, generalPractitioner, link, meta, text, contained []byte
	var extension, modifierExtension []byte
	var managingOrganization []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&patient.ID,
		&identifier,
		&patient.Active,
		&name,
		&telecom,
		&patient.Gender,
		&patient.BirthDate,
		&patient.DeceasedBoolean,
		&patient.DeceasedDateTime,
		&address,
		&maritalStatus,
		&patient.MultipleBirthBoolean,
		&patient.MultipleBirthInteger,
		&photo,
		&contact,
		&communication,
		&generalPractitioner,
		&managingOrganization,
		&link,
		&meta,
		&patient.ImplicitRules,
		&patient.Language,
		&text,
		&contained,
		&extension,
		&modifierExtension,
		&patient.CreatedAt,
		&patient.UpdatedAt,
		&patient.Version,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("patient not found")
		}
		return nil, fmt.Errorf("failed to get patient: %w", err)
	}

	// Unmarshal JSON fields
	if err := unmarshalJSONFields(patient, identifier, name, telecom, address, maritalStatus,
		photo, contact, communication, generalPractitioner, managingOrganization, link,
		meta, text, contained, extension, modifierExtension); err != nil {
		return nil, err
	}

	return patient, nil
}

func (r *PatientRepository) Update(ctx context.Context, patient *models.Patient) error {
	// First get the old values for audit
	oldPatient, err := r.GetByID(ctx, patient.ID)
	if err != nil {
		return err
	}

	query := `
		UPDATE patients SET
			identifier = $2, active = $3, name = $4, telecom = $5, gender = $6,
			birth_date = $7, deceased_boolean = $8, deceased_date_time = $9,
			address = $10, marital_status = $11, multiple_birth_boolean = $12,
			multiple_birth_integer = $13, photo = $14, contact = $15,
			communication = $16, general_practitioner = $17, managing_organization = $18,
			link = $19, meta = $20, implicit_rules = $21, language = $22,
			text = $23, contained = $24, extension = $25, modifier_extension = $26
		WHERE id = $1
		RETURNING updated_at, version
	`

	err = r.db.QueryRowContext(ctx, query,
		patient.ID,
		toJSON(patient.Identifier),
		patient.Active,
		toJSON(patient.Name),
		toJSON(patient.Telecom),
		patient.Gender,
		patient.BirthDate,
		patient.DeceasedBoolean,
		patient.DeceasedDateTime,
		toJSON(patient.Address),
		toJSON(patient.MaritalStatus),
		patient.MultipleBirthBoolean,
		patient.MultipleBirthInteger,
		toJSON(patient.Photo),
		toJSON(patient.Contact),
		toJSON(patient.Communication),
		toJSON(patient.GeneralPractitioner),
		toJSON(patient.ManagingOrganization),
		toJSON(patient.Link),
		toJSON(patient.Meta),
		patient.ImplicitRules,
		patient.Language,
		toJSON(patient.Text),
		toJSON(patient.Contained),
		toJSON(patient.Extension),
		toJSON(patient.ModifierExtension),
	).Scan(&patient.UpdatedAt, &patient.Version)

	if err != nil {
		return fmt.Errorf("failed to update patient: %w", err)
	}

	// Log audit trail
	auditLog := &AuditLog{
		ResourceType: "Patient",
		ResourceID:   patient.ID,
		Action:       "UPDATE",
		OldValues:    mustMarshalJSON(oldPatient),
		NewValues:    mustMarshalJSON(patient),
	}
	
	if err := r.LogAudit(ctx, auditLog); err != nil {
		fmt.Printf("Failed to log audit: %v\n", err)
	}

	return nil
}

func (r *PatientRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Get the patient for audit log
	patient, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	query := `DELETE FROM patients WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete patient: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("patient not found")
	}

	// Log audit trail
	auditLog := &AuditLog{
		ResourceType: "Patient",
		ResourceID:   id,
		Action:       "DELETE",
		OldValues:    mustMarshalJSON(patient),
	}
	
	if err := r.LogAudit(ctx, auditLog); err != nil {
		fmt.Printf("Failed to log audit: %v\n", err)
	}

	return nil
}

func (r *PatientRepository) List(ctx context.Context, params PaginationParams) ([]*models.Patient, PaginationResult, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM patients`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, PaginationResult{}, fmt.Errorf("failed to get patient count: %w", err)
	}

	// Get patients with pagination
	query := `
		SELECT id, identifier, active, name, telecom, gender, birth_date,
			   deceased_boolean, deceased_date_time, address, marital_status,
			   multiple_birth_boolean, multiple_birth_integer, photo, contact,
			   communication, general_practitioner, managing_organization, link,
			   meta, implicit_rules, language, text, contained, extension, 
			   modifier_extension, created_at, updated_at, version
		FROM patients 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, params.Limit, params.Offset)
	if err != nil {
		return nil, PaginationResult{}, fmt.Errorf("failed to list patients: %w", err)
	}
	defer rows.Close()

	var patients []*models.Patient
	for rows.Next() {
		patient := &models.Patient{}
		var identifier, name, telecom, address, maritalStatus, photo, contact []byte
		var communication, generalPractitioner, link, meta, text, contained []byte
		var extension, modifierExtension []byte
		var managingOrganization []byte

		err := rows.Scan(
			&patient.ID,
			&identifier,
			&patient.Active,
			&name,
			&telecom,
			&patient.Gender,
			&patient.BirthDate,
			&patient.DeceasedBoolean,
			&patient.DeceasedDateTime,
			&address,
			&maritalStatus,
			&patient.MultipleBirthBoolean,
			&patient.MultipleBirthInteger,
			&photo,
			&contact,
			&communication,
			&generalPractitioner,
			&managingOrganization,
			&link,
			&meta,
			&patient.ImplicitRules,
			&patient.Language,
			&text,
			&contained,
			&extension,
			&modifierExtension,
			&patient.CreatedAt,
			&patient.UpdatedAt,
			&patient.Version,
		)

		if err != nil {
			return nil, PaginationResult{}, fmt.Errorf("failed to scan patient: %w", err)
		}

		// Unmarshal JSON fields
		if err := unmarshalJSONFields(patient, identifier, name, telecom, address, maritalStatus,
			photo, contact, communication, generalPractitioner, managingOrganization, link,
			meta, text, contained, extension, modifierExtension); err != nil {
			return nil, PaginationResult{}, err
		}

		patients = append(patients, patient)
	}

	if err := rows.Err(); err != nil {
		return nil, PaginationResult{}, fmt.Errorf("failed to iterate patients: %w", err)
	}

	pagination := GetPaginationResult(total, params)
	return patients, pagination, nil
}

// Helper functions
func toJSON(v interface{}) []byte {
	if v == nil {
		return []byte("null")
	}
	data, _ := json.Marshal(v)
	return data
}

func mustMarshalJSON(v interface{}) json.RawMessage {
	data, _ := json.Marshal(v)
	return data
}

func unmarshalJSONFields(patient *models.Patient, fields ...[]byte) error {
	// This would unmarshal all the JSON fields - implementation depends on the models
	// For now, we'll leave this as a placeholder
	return nil
}
