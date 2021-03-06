package gh

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type SearchOptions struct {
	Query       string
	Interactive bool
	SearchIn    string
	Limit       int
	Topic       string
	Language    string
}

func ghCommand() *cobra.Command {
	opts := &SearchOptions{}
	cmd := &cobra.Command{
		Use:   "gh search <query>",
		Short: "search repositories",
		Long: `Search for GitHub repositories.

Filter with the --topic or --lang flag, 
or write a custom query with the -q flag.`,
		Example: `# cli repos with hacktoberfest topic
$ gh search cli --topic=hacktoberfest

# custom search with GitHub syntax
$ gh search -q="org:cli created:>2019-01-01"`,
		Args:          cobra.MaximumNArgs(1),
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Query == "" {
				if len(args) < 1 {
					return errors.New("Error: empty query")
				}
				opts.Query = args[0]
				if opts.SearchIn != "" {
					searchIn := strings.ToLower(opts.SearchIn)
					if searchIn != "name" && searchIn != "description" && searchIn != "readme" {
						return errors.New(`--in argument must be "name", "description", or "readme"`)
					}
				}
				if opts.Limit <= 0 {
					return errors.New("invalid limit")
				}
				opts.Query = prepareQuery(opts)
			}

			return runSearch(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Topic, "topic", "t", "", `Specify a topic`)
	cmd.Flags().StringVarP(&opts.SearchIn, "in", "i", "", `Search in "name", "description", or "readme"`)
	cmd.Flags().StringVarP(&opts.Language, "lang", "l", "", "Search by programming language")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 50, "Max number of search results")
	cmd.Flags().StringVarP(&opts.Query, "query", "q", "", "Query in GitHub syntax")
	return cmd
}

func runSearch(opts *SearchOptions) error {
	results, total, err := searchRepos(opts)
	if err != nil {
		return err
	}

	if total == 0 {
		fmt.Printf(`No results found for "%s"%s`, opts.Query, "\n")
	}

	var repoStrs []string
	for i, repo := range results {
		repoStrs = append(repoStrs, prettyPrint(i+1, &repo))
	}

	numResults := len(repoStrs)

	selector := &survey.Select{
		Message:  fmt.Sprintf("%d/%d Results\n", numResults, total),
		Options:  repoStrs,
		PageSize: 10,
	}

	var selection string
	err = survey.AskOne(selector, &selection,
		survey.WithIcons(func(icons *survey.IconSet) {
			icons.Question.Text = "::"
			icons.Question.Format = "yellow+hb"
		}),
	)
	if err != nil {
		return nil
	}

	n, err := strconv.Atoi(strings.Split(selection, " ")[0])
	if err != nil {
		return err
	}
	selectedRepo := results[n-1]
	args := []string{"repo", "view", selectedRepo.NameWithOwner}
	stdOut, _, err := gh.Exec(args...)
	if err != nil {
		return err
	}
	fmt.Print(stdOut.String())

	return nil
}

func prettyPrint(i int, repo *Repository) string {
	out := fmt.Sprintf("%d %s\n", i, repo.NameWithOwner)

	dscript := repo.Description
	if len(dscript) > 100 {
		dscript = dscript[0:97]
		dscript += "..."
	}
	out += fmt.Sprintf("\t%s\n", dscript)

	lang := repo.PrimaryLanguage.Name

	if lang != "" {
		out += fmt.Sprintf("\tLanguage: %s\n", lang)
	}

	if repo.StargazerCount >= 1000 {
		out += fmt.Sprintf("\t??? %.1fk", float32(repo.StargazerCount)/1000.0)
	} else {
		out += fmt.Sprintf("\t??? %d", repo.StargazerCount)
	}
	return out
}

//func main() {
//	cmd := rootCmd()
//	if err := cmd.Execute(); err != nil {
//		fmt.Fprintf(os.Stderr, "%s\n", err)
//		os.Exit(1)
//	}
//}

func init() {
	rootCmd.AddCommand(ghCommand)
}
