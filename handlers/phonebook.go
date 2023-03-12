package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)


type Phonebook struct {
	l *log.Logger
}

type Contact struct {
	ID uuid.UUID `json:"id"`
	Name string `json:"name"`
	Number string `json:"number"`
	Email string `json:"email"`
}

type Contacts []*Contact
// TODO: Add proper log handler and error log handler
func NewPhonebook(l *log.Logger) *Phonebook {
	return &Phonebook{l}
}

func (p *Phonebook) ListContact() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p.l.Printf("Calling list contacts\n")
		e := json.NewEncoder(w)
		if err:= e.Encode(contactList); err != nil {
			http.Error(w, "Internal error trying to encode list", http.StatusInternalServerError)
		}
	})
}

func (p *Phonebook) AddContact() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p.l.Printf("Calling add contacts\n")
		c := &Contact{}
		e := json.NewDecoder(r.Body)
		if err := e.Decode(c); err != nil {
			http.Error(w, "Unable to read body", http.StatusBadRequest)
		}
		
		c.ID = uuid.New()
		contactList = append(contactList, c)

		w.Write([]byte(c.ID.String()))
	})
}

func (p *Phonebook) DeleteContact() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p.l.Printf("Calling delete contacts")
		id, err := uuid.Parse(strings.TrimPrefix(r.URL.Path, "/delete/"))
		if err != nil {
			http.Error(w, "Unable to get id", http.StatusBadRequest)
		}
		p.l.Printf(" id=%s\n", id)


		_, i, err := findContact(id)
		if err == ErrContactNotFound {
			http.Error(w, "Contact not found", http.StatusNotFound)
		}

		if err != nil {
			http.Error(w, "Internal error", http.StatusBadRequest)
		}

		contactList = append(contactList[:i], contactList[i+1:]...)

    w.WriteHeader(http.StatusOK)

	})
}

func (p *Phonebook) UpdateContact() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p.l.Printf("Calling update contacts")
		id, err := uuid.Parse(strings.TrimPrefix(r.URL.Path, "/update/"))
		if err != nil {
			http.Error(w, "Unable to get id", http.StatusBadRequest)
		}

		newC := &Contact{}
		e := json.NewDecoder(r.Body)
		if err := e.Decode(newC); err != nil {
			http.Error(w, "Unable to read body", http.StatusBadRequest)
		}

		_, pos, err := findContact(id)
		if err == ErrContactNotFound {
			http.Error(w, "Contact not found", http.StatusNotFound)
		}

		if err != nil {
			http.Error(w, "Internal error", http.StatusBadRequest)
		}

		newC.ID = id
		contactList[pos] = newC

    w.WriteHeader(http.StatusOK)
	})
}

func (p *Phonebook) FindContact() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p.l.Printf("Calling find contacts")
		id, err := uuid.Parse(strings.TrimPrefix(r.URL.Path, "/find/"))
		if err != nil {
			http.Error(w, "Unable to get id", http.StatusBadRequest)
		}

		c, _, err := findContact(id)
		if err == ErrContactNotFound {
			http.Error(w, "Contact not found", http.StatusNotFound)
		}

		if err != nil {
			http.Error(w, "Internal error", http.StatusBadRequest)
		}


		e := json.NewEncoder(w)
		if err:= e.Encode(c); err != nil {
			http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		}
	})
}

func (p *Phonebook) FindContactByName() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p.l.Printf("Calling find contacts by name")
		name := strings.TrimPrefix(r.URL.Path, "/find-by-name/")
		
		c, err := findByName(name)

		if err == ErrContactNotFound {
			http.Error(w, "Contact not found", http.StatusNotFound)
		}

		if err != nil {
			http.Error(w, "Unable to find contact", http.StatusNotFound)
		}

		e := json.NewEncoder(w)
		if err:= e.Encode(c); err != nil {
			http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		}
	})
}

var ErrContactNotFound = fmt.Errorf("Contact not found")

func findContact(id uuid.UUID) (*Contact, int, error) {
	for i, c := range contactList {
		if c.ID == id {
			return c, i, nil
		}
	}

	return nil, -1, ErrContactNotFound
} 

func findByName(name string) (*Contact, error) {
	for _, c := range contactList {
		if strings.Contains(c.Name, name) {
			return c, nil
		}
	}

	return nil, ErrContactNotFound
}

var contactList = Contacts{
	&Contact{
		ID:        uuid.New(),
		Name:      "Juju",
		Number:     "48 99119-9999",
		Email:     "juju@gmail.com",
	},
	&Contact{
		ID:        uuid.New(),
		Name:      "Tutu",
		Number:     "48 99119-9999",
		Email:     "tutu@gmail.com",
	},
}
