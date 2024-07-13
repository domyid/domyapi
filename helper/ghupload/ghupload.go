package ghupload

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/google/go-github/v59/github"

	"golang.org/x/oauth2"
)

func GithubListFiles(GitHubAccessToken, githubOrg, githubRepo, path string) ([]*github.RepositoryContent, error) {
	// Konfigurasi koneksi ke GitHub menggunakan token akses
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GitHubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Mendapatkan daftar file dari repositori
	_, directoryContent, _, err := client.Repositories.GetContents(ctx, githubOrg, githubRepo, path, nil)
	if err != nil {
		return nil, err
	}

	// Tambahkan logging untuk melihat data yang dikembalikan
	fmt.Printf("GithubListFiles: %v\n", directoryContent)

	return directoryContent, nil
}

func GithubUpload(accessToken, authorName, authorEmail string, fileHeader *multipart.FileHeader, org, repo, path string, replace bool, fileContent io.Reader) (*github.RepositoryContentResponse, *github.Response, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	content, err := io.ReadAll(fileContent)
	if err != nil {
		return nil, nil, err
	}

	opts := &github.RepositoryContentFileOptions{
		Message: github.String("Upload file"),
		Content: content,
		Branch:  github.String("main"),
		Author: &github.CommitAuthor{
			Name:  github.String(authorName),
			Email: github.String(authorEmail),
		},
	}

	fileContentSHA := ""
	if replace {
		currentContent, _, _, err := client.Repositories.GetContents(ctx, org, repo, path, nil)
		if err == nil {
			fileContentSHA = currentContent.GetSHA()
			opts.SHA = github.String(fileContentSHA)
		}
	}

	return client.Repositories.CreateFile(ctx, org, repo, path, opts)
}

func GithubUpdateFile(GitHubAccessToken, GitHubAuthorName, GitHubAuthorEmail string, fileHeader *multipart.FileHeader, githubOrg, githubRepo, pathFile string) (*github.RepositoryContentResponse, *github.Response, error) {
	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	// Read the file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}

	// Konfigurasi koneksi ke GitHub menggunakan token akses
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GitHubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Get the current file to retrieve the SHA
	currentContent, _, _, err := client.Repositories.GetContents(ctx, githubOrg, githubRepo, pathFile, nil)
	if err != nil {
		return nil, nil, err
	}

	opts := &github.RepositoryContentFileOptions{
		Message: github.String("Update file"),
		Content: fileContent,
		Branch:  github.String("main"),
		SHA:     github.String(currentContent.GetSHA()),
		Author: &github.CommitAuthor{
			Name:  github.String(GitHubAuthorName),
			Email: github.String(GitHubAuthorEmail),
		},
	}

	// Update the file in the repository
	return client.Repositories.UpdateFile(ctx, githubOrg, githubRepo, pathFile, opts)
}

func GithubDeleteFile(GitHubAccessToken, GitHubAuthorName, GitHubAuthorEmail, githubOrg, githubRepo, pathFile string) (*github.RepositoryContentResponse, *github.Response, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GitHubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Get the current file to retrieve the SHA
	currentContent, _, _, err := client.Repositories.GetContents(ctx, githubOrg, githubRepo, pathFile, nil)
	if err != nil {
		return nil, nil, err
	}

	opts := &github.RepositoryContentFileOptions{
		Message: github.String("Delete file"),
		Branch:  github.String("main"),
		SHA:     github.String(currentContent.GetSHA()),
		Author: &github.CommitAuthor{
			Name:  github.String(GitHubAuthorName),
			Email: github.String(GitHubAuthorEmail),
		},
	}

	// Delete the file from the repository
	deleteResponse, response, err := client.Repositories.DeleteFile(ctx, githubOrg, githubRepo, pathFile, opts)
	if err != nil {
		return nil, response, err
	}

	return deleteResponse, response, nil
}
