package client

import (
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var cpfDigits = regexp.MustCompile(`\d`)

type Client struct {
	ID        uuid.UUID
	Name      string
	CPF       string
	Email     string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func NewClient(name, cpf, email, phone string) (*Client, error) {
	if len(name) < 3 {
		return nil, ErrInvalidName
	}

	cpf = normalizeCPF(cpf)
	if len(cpf) != 11 {
		return nil, ErrInvalidCPF
	}

	if email != "" {
		if _, err := mail.ParseAddress(email); err != nil {
			return nil, ErrInvalidEmail
		}
	}

	now := time.Now()

	return &Client{
		ID:        uuid.New(),
		Name:      name,
		CPF:       cpf,
		Email:     email,
		Phone:     phone,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (c *Client) SoftDelete() {
	now := time.Now()
	c.DeletedAt = &now
}

func (c *Client) Update(name, email, phone string) error {
	if len(name) < 3 {
		return ErrInvalidName
	}
	if email != "" {
		if _, err := mail.ParseAddress(email); err != nil {
			return ErrInvalidEmail
		}
	}
	c.Name = name
	c.Email = email
	c.Phone = phone
	c.UpdatedAt = time.Now()
	return nil
}

func normalizeCPF(cpf string) string {
	return strings.Join(cpfDigits.FindAllString(cpf, -1), "")
}
