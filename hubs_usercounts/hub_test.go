package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	cdp "github.com/chromedp/chromedp"
)

func TestHubs(t *testing.T) {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := cdp.New(ctxt, cdp.WithLog(log.Printf))
	if err != nil {
		t.Fatal(err)
	}

	// Login to hub
	err = loginToHub(ctxt, c, `ekrchan.sit1+1@gmail.com`, `******`)
	if err != nil {
		t.Error(err)
	}

	// Check hub counts
	res, err := checkUserCounts(ctxt, c)
	if err != nil {
		t.Error(err)
	}

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		t.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Success: %s", res)
}

func loginToHub(ctxt context.Context, c *cdp.CDP, email string, pw string) error {
	// Force max timeout of 15 seconds for retrieving and processing the data
	var cancel func()
	ctxt, cancel = context.WithTimeout(ctxt, 15*time.Second)
	defer cancel()

	const hubUrl = `https://sit1.pressly.com/s3upgrade03/englishb-en-ca`
	const emailSel = `//input[@name="email"]`
	const pwSel = `//input[@name="password"]`
	const submitButton = `button[type="submit"]`

	// Login to hub
	if err := c.Run(ctxt, cdp.Tasks{
		cdp.Navigate(hubUrl),
		cdp.WaitVisible(emailSel),
		cdp.SendKeys(emailSel, email),
		cdp.SendKeys(pwSel, pw),
		cdp.Click(submitButton),
		cdp.WaitNotPresent(emailSel)}); err != nil {
		return fmt.Errorf("Could not navigate to the hub: %v", err)
	}

	return nil
}

func checkUserCounts(ctxt context.Context, c *cdp.CDP) (string, error) {
	// Force max timeout of 15 seconds for retrieving and processing the data
	var cancel func()
	ctxt, cancel = context.WithTimeout(ctxt, 15*time.Second)
	defer cancel()

	const submitButton = `button[type="submit"]`
	const coverCountSel = `.px-member-count`
	const collabCountSel = `a[class="px-team-header-active"] > strong:nth-child(2)`
	const memberCountSel = `.zt7ggh-2 > li:nth-child(2) > a:nth-child(1) > strong:nth-child(2)`

	var hubCoverString string
	var collaboratorString string
	var memberString string
	var totalCount int
	var collabCount int
	var memberCount int

	// Grab count on hub cover
	if err := c.Run(ctxt, cdp.Text(coverCountSel, &hubCoverString)); err != nil {
		return "", fmt.Errorf("Could not retrieve count from hub cover: %v", err)
	}

	if total, err := strconv.Atoi(hubCoverString); err != nil {
		return "", fmt.Errorf("Could not convert hub cover count to int: %v", err)
	} else {
		totalCount = total
	}

	// Navigate to team page and get counts
	if err := c.Run(ctxt, cdp.Tasks{
		cdp.Click(coverCountSel),
		cdp.WaitNotPresent(coverCountSel),
		cdp.Text(collabCountSel, &collaboratorString),
		cdp.Text(memberCountSel, &memberString),
		cdp.Sleep(5 * time.Second)}); err != nil {
		return "", fmt.Errorf("Could not retrieve counts from team page: %v", err)
	}

	if collab, err := strconv.Atoi(collaboratorString); err != nil {
		return "", fmt.Errorf("Could not convert collaborator count to int: %v", err)
	} else {
		collabCount = collab
	}

	if member, err := strconv.Atoi(memberString); err != nil {
		return "", fmt.Errorf("Could not convert member count to int: %v", err)
	} else {
		memberCount = member
	}

	// Check collaborators + members = hub cover count
	if totalCount != collabCount+memberCount {
		err := fmt.Errorf("Total count does not match member count + collaborator count. Total Count: %v. Collaborator Count: %v. Member Count: %v", totalCount, collabCount, memberCount)
		return "", err
	}

	return fmt.Sprintf("Total Count: %v. Collaborator Count: %v. Member Count: %v", totalCount, collabCount, memberCount), nil
}
