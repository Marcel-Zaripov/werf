package main

import (
	"encoding/json"
	"fmt"

	"github.com/flant/dapp/pkg/deploy"
	"github.com/flant/dapp/pkg/ruby2go"
	"github.com/flant/dapp/pkg/secret"
)

func main() {
	ruby2go.RunCli("deploy", func(args map[string]interface{}) (interface{}, error) {
		cmd, err := ruby2go.CommandFieldFromArgs(args)
		if err != nil {
			return nil, err
		}

		switch cmd {
		case "secret_key_generate":
			key, err := secret.GenerateAexSecretKey()
			if err != nil {
				return nil, err
			}

			fmt.Printf("DAPP_SECRET_KEY=%s\n", string(key))

			return nil, nil
		case "secret_generate", "secret_extract":
			projectDir, err := ruby2go.StringOptionFromArgs("project_dir", args)
			if err != nil {
				return nil, err
			}

			rawOptions, err := ruby2go.StringOptionFromArgs("raw_command_options", args)
			if err != nil {
				return nil, err
			}

			options := &secretGenerateOptions{}
			err = json.Unmarshal([]byte(rawOptions), options)
			if err != nil {
				return nil, err
			}

			s, err := deploy.GetSecret(projectDir)
			if err != nil {
				return nil, err
			}

			var secretGenerator *deploy.SecretGenerator
			switch cmd {
			case "secret_generate":
				if secretGenerator, err = newSecretGenerateGenerator(s); err != nil {
					return nil, err
				}

				return nil, secretGenerate(secretGenerator, *options)
			case "secret_extract":
				if secretGenerator, err = newSecretExtractGenerator(s); err != nil {
					return nil, err
				}

				return nil, secretExtract(secretGenerator, *options)
			}
		case "secret_regenerate":
			projectDir, err := ruby2go.StringOptionFromArgs("project_dir", args)
			if err != nil {
				return nil, err
			}

			oldKey, err := ruby2go.StringOptionFromArgs("old_key", args)
			if err != nil {
				return nil, err
			}

			secretValuesPaths, err := ruby2go.StringArrayOptionFromArgs("secret_values_paths", args)
			if err != nil {
				return nil, err
			}

			newSecret, err := deploy.GetSecret(projectDir)
			if err != nil {
				return nil, err
			}

			oldSecret, err := secret.NewSecret([]byte(oldKey))
			if err != nil {
				return nil, err
			}

			return nil, SecretsRegenerate(newSecret, oldSecret, projectDir, secretValuesPaths...)
		case "secret_edit":
		default:
			return nil, fmt.Errorf("command `%s` isn't supported", cmd)
		}

		return nil, nil
	})
}
