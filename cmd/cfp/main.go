package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/core-file-privacy/core-file-privacy/internal/app"
	"github.com/core-file-privacy/core-file-privacy/internal/files"
	"github.com/core-file-privacy/core-file-privacy/internal/prompt"
	"github.com/core-file-privacy/core-file-privacy/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cfp",
	Short: "Encrypt, hide and restore private files from your terminal",
	Long: `Core File Privacy - Encrypt, hide and restore private files from your terminal.

A CLI tool for encrypting, hiding, verifying, decrypting and changing
passwords of files or folders using a user-provided password or a
secure auto-generated key.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		checkInstallation(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var hideCmd = &cobra.Command{
	Use:   "hide <path>",
	Short: "Encrypt and optionally hide a file or folder",
	Long: `Encrypt a file or folder into a .cfp container.

If the path is a directory, you must use --archive to create a single
encrypted container.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		generateKey, _ := cmd.Flags().GetBool("generate-key")
		keep, _ := cmd.Flags().GetBool("keep")
		deleteOriginal, _ := cmd.Flags().GetBool("delete-original")
		hidden, _ := cmd.Flags().GetBool("hidden")
		output, _ := cmd.Flags().GetString("output")
		archive, _ := cmd.Flags().GetBool("archive")
		nameMode, _ := cmd.Flags().GetString("name-mode")
		profile, _ := cmd.Flags().GetString("profile")
		yes, _ := cmd.Flags().GetBool("yes")

		var password string
		if !generateKey {
			var err error
			password, err = prompt.ReadPasswordTwice()
			if err != nil {
				return err
			}
		}

		result, err := app.Hide(app.HideOptions{
			Path:           args[0],
			Password:       password,
			GenerateKey:    generateKey,
			Keep:           keep,
			DeleteOriginal: deleteOriginal,
			Hidden:         hidden,
			Output:         output,
			Archive:        archive,
			NameMode:       nameMode,
			Profile:        profile,
			Yes:            yes,
		})
		if err != nil {
			return err
		}

		if result.GeneratedKey != "" {
			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, "Generated key:")
			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, result.GeneratedKey)
			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, "Save this key in a password manager.")
			fmt.Fprintln(os.Stderr, "Anyone with this key can decrypt the file.")
			fmt.Fprintln(os.Stderr, "CFP cannot recover it if lost.")
			fmt.Fprintln(os.Stderr)
			prompt.WaitForEnter("Press Enter once you have saved it.")
		}

		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Encrypted:")
		fmt.Fprintf(os.Stderr, "  %s\n", result.OutputPath)
		fmt.Fprintln(os.Stderr)
		if result.OriginalKept {
			fmt.Fprintln(os.Stderr, "Original:")
			fmt.Fprintln(os.Stderr, "  kept")
		} else {
			fmt.Fprintln(os.Stderr, "Original:")
			fmt.Fprintln(os.Stderr, "  deleted")
		}
		if hidden {
			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, "Hidden:")
			fmt.Fprintln(os.Stderr, "  yes")
		}

		return nil
	},
}

var showCmd = &cobra.Command{
	Use:   "show <file.cfp>",
	Short: "Decrypt a .cfp container",
	Long:  `Decrypt a .cfp container and restore the original file or folder.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")
		keep, _ := cmd.Flags().GetBool("keep")
		deleteEncrypted, _ := cmd.Flags().GetBool("delete-encrypted")
		force, _ := cmd.Flags().GetBool("force")
		yes, _ := cmd.Flags().GetBool("yes")

		password, err := prompt.ReadPassword("Enter password/key: ")
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr)

		result, err := app.Show(app.ShowOptions{
			Path:            args[0],
			Password:        password,
			Output:          output,
			Keep:            keep,
			DeleteEncrypted: deleteEncrypted,
			Force:           force,
			Yes:             yes,
		})
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stderr, "Decrypted:")
		fmt.Fprintf(os.Stderr, "  %s\n", result.OutputPath)

		if deleteEncrypted {
			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, "Encrypted file:")
			fmt.Fprintln(os.Stderr, "  deleted")
		}

		return nil
	},
}

var verifyCmd = &cobra.Command{
	Use:   "verify <file.cfp>",
	Short: "Verify container integrity without decrypting",
	Long:  `Verify that a password/key is correct and the container is intact.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		password, err := prompt.ReadPassword("Enter password/key: ")
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr)

		err = app.Verify(app.VerifyOptions{
			Path:     args[0],
			Password: password,
		})
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stderr, "Password/key is valid.")
		fmt.Fprintln(os.Stderr, "Container integrity verified.")

		return nil
	},
}

var infoCmd = &cobra.Command{
	Use:   "info <file.cfp>",
	Short: "Show container information",
	Long:  `Display public information from the container header. Use --unlock to show private metadata.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		unlock, _ := cmd.Flags().GetBool("unlock")

		var password string
		if unlock {
			var err error
			password, err = prompt.ReadPassword("Enter password/key: ")
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stderr)
		}

		result, err := app.Info(app.InfoOptions{
			Path:     args[0],
			Password: password,
			Unlock:   unlock,
		})
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stderr, "Core File Privacy container")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Version: CFP%d\n", result.Header.Version)
		fmt.Fprintf(os.Stderr, "Type: %s\n", result.Header.Type)
		fmt.Fprintf(os.Stderr, "Cipher: %s\n", result.Header.Cipher)
		fmt.Fprintf(os.Stderr, "KDF: %s\n", result.Header.KDF)
		fmt.Fprintf(os.Stderr, "Profile: %s\n", result.Header.Profile)
		fmt.Fprintf(os.Stderr, "Created at: %s\n", result.Header.CreatedAt)
		fmt.Fprintf(os.Stderr, "Metadata encrypted: %v\n", result.Header.MetadataEncrypted)

		if result.Metadata != nil {
			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, "Original name:", result.Metadata.OriginalName)
			fmt.Fprintln(os.Stderr, "Original type:", result.Metadata.OriginalType)
			if result.Metadata.OriginalMode != "" {
				fmt.Fprintln(os.Stderr, "Original mode:", result.Metadata.OriginalMode)
			}
			if result.Metadata.OriginalSize > 0 {
				fmt.Fprintln(os.Stderr, "Original size:", result.Metadata.OriginalSize)
			}
			fmt.Fprintln(os.Stderr, "Archived:", result.Metadata.Archived)
		}

		return nil
	},
}

var rekeyCmd = &cobra.Command{
	Use:   "rekey <file.cfp>",
	Short: "Change the password/key of a container",
	Long:  `Decrypt and re-encrypt the container with a new password or key.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		generateKey, _ := cmd.Flags().GetBool("generate-key")
		profile, _ := cmd.Flags().GetString("profile")
		yes, _ := cmd.Flags().GetBool("yes")

		fmt.Fprintln(os.Stderr, "Enter current password/key:")
		oldPassword, err := prompt.ReadPassword("Password: ")
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr)

		var newPassword string
		if !generateKey {
			fmt.Fprintln(os.Stderr, "Enter new password/key:")
			newPassword, err = prompt.ReadPasswordTwice()
			if err != nil {
				return err
			}
		}

		result, err := app.Rekey(app.RekeyOptions{
			Path:        args[0],
			OldPassword: oldPassword,
			NewPassword: newPassword,
			GenerateKey: generateKey,
			Profile:     profile,
			Yes:         yes,
		})
		if err != nil {
			return err
		}

		if result.GeneratedKey != "" {
			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, "Generated new key:")
			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, result.GeneratedKey)
			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, "Save this key in a password manager.")
			fmt.Fprintln(os.Stderr, "Anyone with this key can decrypt the file.")
			fmt.Fprintln(os.Stderr, "CFP cannot recover it if lost.")
			fmt.Fprintln(os.Stderr)
			prompt.WaitForEnter("Press Enter once you have saved it.")
		}

		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Rekey completed successfully.")

		return nil
	},
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install cfp to user tools directory",
	Long: `Install cfp to the user tools directory and add it to PATH.

Default installation directory:
  Linux/macOS: ~/core-file-privacy/
  Windows: %USERPROFILE%\core-file-privacy\`,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		noPath, _ := cmd.Flags().GetBool("no-path")
		pathOnly, _ := cmd.Flags().GetBool("path-only")
		target, _ := cmd.Flags().GetString("target")
		yes, _ := cmd.Flags().GetBool("yes")

		return app.Install(app.InstallOptions{
			Force:    force,
			NoPath:   noPath,
			PathOnly: pathOnly,
			Target:   target,
			Yes:      yes,
		})
	},
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove cfp installation",
	Long:  `Remove cfp from the user tools directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		yes, _ := cmd.Flags().GetBool("yes")

		return app.Uninstall(app.UninstallOptions{
			Yes: yes,
		})
	},
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system status and configuration",
	Long:  `Diagnose the system configuration and verify installation status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.Doctor()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.GetFullVersion())
	},
}

func init() {
	hideCmd.Flags().BoolP("generate-key", "g", false, "Generate a secure key instead of usar una contraseña")
	hideCmd.Flags().BoolP("keep", "k", false, "Keep the original file after encryption")
	hideCmd.Flags().BoolP("delete-original", "d", false, "Delete the original file after encryption")
	hideCmd.Flags().BoolP("hidden", "H", false, "Hide the encrypted file")
	hideCmd.Flags().StringP("output", "o", "", "Output file path")
	hideCmd.Flags().BoolP("archive", "a", false, "Archive directory before encryption")
	hideCmd.Flags().String("name-mode", "keep", "Name mode: keep or random")
	hideCmd.Flags().String("profile", "default", "Security profile: fast, default, or paranoid")
	hideCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompts")

	showCmd.Flags().StringP("output", "o", "", "Output file path")
	showCmd.Flags().BoolP("keep", "k", false, "Keep the encrypted file after decryption")
	showCmd.Flags().BoolP("delete-encrypted", "d", false, "Delete the encrypted file after decryption")
	showCmd.Flags().BoolP("force", "f", false, "Overwrite existing files")
	showCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompts")

	infoCmd.Flags().Bool("unlock", false, "Decrypt and show private metadata")

	rekeyCmd.Flags().BoolP("generate-key", "g", false, "Generate a new secure key")
	rekeyCmd.Flags().String("profile", "", "Security profile: fast, default, or paranoid")
	rekeyCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompts")

	installCmd.Flags().Bool("force", false, "Force reinstall even if already installed")
	installCmd.Flags().Bool("no-path", false, "Copy binary but don't modify PATH")
	installCmd.Flags().Bool("path-only", false, "Only update PATH, don't copy binary")
	installCmd.Flags().String("target", "", "Custom installation target path")
	installCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompts")

	uninstallCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompts")

	rootCmd.AddCommand(hideCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(verifyCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(rekeyCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(versionCmd)

	rootCmd.Version = version.GetVersion()
	rootCmd.SetVersionTemplate("cfp version {{.Version}}\n")
}

func checkInstallation(cmd *cobra.Command) {
	if cmd.Name() == "install" || cmd.Name() == "uninstall" || cmd.Name() == "doctor" || cmd.Name() == "version" || cmd.Name() == "help" || cmd.Name() == "completion" {
		return
	}

	currentBinary, err := os.Executable()
	if err != nil {
		return
	}

	currentBinary, err = filepath.EvalSymlinks(currentBinary)
	if err != nil {
		return
	}

	expectedBinary := files.GetInstallBinaryPath()

	if currentBinary != expectedBinary && !files.FileExists(expectedBinary) {
		fmt.Fprintln(os.Stderr, "Core File Privacy is not installed in your user tools directory.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Install it now?\nTarget: %s\n\n", expectedBinary)
		fmt.Fprint(os.Stderr, "[Y/n] ")

		var response string
		fmt.Scanln(&response)

		if response == "" || response == "y" || response == "Y" || response == "yes" || response == "Yes" {
			err := app.Install(app.InstallOptions{Yes: true})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Installation failed: %v\n", err)
			}
			os.Exit(0)
		}
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
