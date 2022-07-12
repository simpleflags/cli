package ui

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
)

type PromptContent struct {
	ErrMessage string
	Label      string
}

func TextInput(pc PromptContent) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.ErrMessage)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     pc.Label + ":",
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	return result
}

func EmailInput(pc PromptContent) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.ErrMessage)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     pc.Label + ":",
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	return result
}

func Password(pc PromptContent) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.ErrMessage)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     "Password:",
		Templates: templates,
		Validate:  validate,
		Mask:      '*',
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	return result
}
