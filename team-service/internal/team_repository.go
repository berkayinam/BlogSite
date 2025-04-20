package internal

import (
	"database/sql"
	"errors"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) CreateTeam(name, description, creatorUsername string) (*Team, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var team Team
	err = tx.QueryRow(`
		INSERT INTO teams (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at
	`, name, description).Scan(&team.ID, &team.Name, &team.Description, &team.CreatedAt, &team.UpdatedAt)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`
		INSERT INTO team_members (team_id, username, role)
		VALUES ($1, $2, 'admin')
	`, team.ID, creatorUsername)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &team, nil
}

func (r *TeamRepository) GetTeam(teamID int) (*Team, error) {
	team := &Team{}
	err := r.db.QueryRow(`
		SELECT id, name, description, created_at, updated_at
		FROM teams WHERE id = $1
	`, teamID).Scan(&team.ID, &team.Name, &team.Description, &team.CreatedAt, &team.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("team not found")
		}
		return nil, err
	}

	rows, err := r.db.Query(`
		SELECT team_id, username, role, joined_at
		FROM team_members WHERE team_id = $1
	`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var member TeamMember
		err := rows.Scan(&member.TeamID, &member.Username, &member.Role, &member.JoinedAt)
		if err != nil {
			return nil, err
		}
		team.Members = append(team.Members, member)
	}

	return team, nil
}

func (r *TeamRepository) GetUserTeams(username string) ([]Team, error) {
	rows, err := r.db.Query(`
		SELECT t.id, t.name, t.description, t.created_at, t.updated_at
		FROM teams t
		INNER JOIN team_members tm ON t.id = tm.team_id
		WHERE tm.username = $1
	`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []Team
	for rows.Next() {
		var team Team
		err := rows.Scan(&team.ID, &team.Name, &team.Description, &team.CreatedAt, &team.UpdatedAt)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	return teams, nil
}

func (r *TeamRepository) UpdateTeam(teamID int, name, description string) error {
	result, err := r.db.Exec(`
		UPDATE teams
		SET name = $1, description = $2
		WHERE id = $3
	`, name, description, teamID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("team not found")
	}

	return nil
}

func (r *TeamRepository) DeleteTeam(teamID int) error {
	result, err := r.db.Exec("DELETE FROM teams WHERE id = $1", teamID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("team not found")
	}

	return nil
}

func (r *TeamRepository) IsMemberAdmin(teamID int, username string) (bool, error) {
	var role string
	err := r.db.QueryRow(`
		SELECT role FROM team_members
		WHERE team_id = $1 AND username = $2
	`, teamID, username).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return role == "admin", nil
}

func (r *TeamRepository) IsMember(teamID int, username string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM team_members
			WHERE team_id = $1 AND username = $2
		)
	`, teamID, username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *TeamRepository) CreateInvite(teamID int, inviterUsername, inviteeUsername string) (*TeamInvite, error) {
	invite := &TeamInvite{}
	err := r.db.QueryRow(`
		INSERT INTO team_invites (team_id, inviter_username, invitee_username, status)
		VALUES ($1, $2, $3, 'pending')
		RETURNING id, team_id, inviter_username, invitee_username, status, created_at, updated_at
	`, teamID, inviterUsername, inviteeUsername).Scan(
		&invite.ID, &invite.TeamID, &invite.InviterUsername,
		&invite.InviteeUsername, &invite.Status, &invite.CreatedAt, &invite.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return invite, nil
}

func (r *TeamRepository) GetInvite(inviteID int) (*TeamInvite, error) {
	invite := &TeamInvite{}
	err := r.db.QueryRow(`
		SELECT id, team_id, inviter_username, invitee_username, status, created_at, updated_at
		FROM team_invites WHERE id = $1
	`, inviteID).Scan(
		&invite.ID, &invite.TeamID, &invite.InviterUsername,
		&invite.InviteeUsername, &invite.Status, &invite.CreatedAt, &invite.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invite not found")
		}
		return nil, err
	}

	return invite, nil
}

func (r *TeamRepository) GetUserInvites(username string) ([]TeamInvite, error) {
	rows, err := r.db.Query(`
		SELECT id, team_id, inviter_username, invitee_username, status, created_at, updated_at
		FROM team_invites
		WHERE invitee_username = $1 AND status = 'pending'
	`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invites []TeamInvite
	for rows.Next() {
		var invite TeamInvite
		err := rows.Scan(
			&invite.ID, &invite.TeamID, &invite.InviterUsername,
			&invite.InviteeUsername, &invite.Status, &invite.CreatedAt, &invite.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		invites = append(invites, invite)
	}

	return invites, nil
}

func (r *TeamRepository) UpdateInviteStatus(inviteID int, status string) error {
	result, err := r.db.Exec(`
		UPDATE team_invites
		SET status = $1
		WHERE id = $2 AND status = 'pending'
	`, status, inviteID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("invite not found or already processed")
	}

	return nil
}

func (r *TeamRepository) AddMember(teamID int, username string) error {
	_, err := r.db.Exec(`
		INSERT INTO team_members (team_id, username, role)
		VALUES ($1, $2, 'member')
	`, teamID, username)
	return err
}

func (r *TeamRepository) RemoveMember(teamID int, username string) error {
	result, err := r.db.Exec(`
		DELETE FROM team_members
		WHERE team_id = $1 AND username = $2 AND role != 'admin'
	`, teamID, username)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("member not found or is an admin")
	}

	return nil
}