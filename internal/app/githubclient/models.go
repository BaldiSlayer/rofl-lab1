package githubclient

type Gist struct {
	Description string
	Files       []GistFile
}

type GistFile struct {
	Name    string
	Content string
}
