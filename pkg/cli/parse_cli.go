package cli

import (
	"reflect"
)

type Parser struct {
	CommandArgs     []string
	acceptedCommand []AcceptedCommand
}

func (p *Parser) Parse() []AcceptedCommand {
	var acceptedCommand []AcceptedCommand
	for i := 0; i < len(p.CommandArgs); i++ {
		if p.CommandArgs[i][:1] == "-" {
			x := AcceptedCommand{}
			x.CommandName = p.CommandArgs[i]
			c := p.CommandArgs[i+1:]
			if !reflect.DeepEqual(c, []string{}) && c[0][:1] != "-" {
				for j := 0; j < len(c); j++ {
					if j == len(c)-1 {
						if c[j][:1] != "-" {
							x.CommandValue = append(x.CommandValue, c[j])
							acceptedCommand = append(acceptedCommand, x)
						} else {
							acceptedCommand = append(acceptedCommand, x)
						}
					} else {
						if c[j][:1] != "-" {
							x.CommandValue = append(x.CommandValue, c[j])
						} else {
							acceptedCommand = append(acceptedCommand, x)
							break
						}
					}
				}
			} else {
				acceptedCommand = append(acceptedCommand, x)
			}
		}
	}
	return acceptedCommand
}
