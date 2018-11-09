package main

import (
	"context"
	"log"
	"time"

	cdp "github.com/chromedp/chromedp"
)

func main() {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := cdp.New(ctxt, cdp.WithLog(log.Printf))
	if err != nil {
		log.Fatal(err)
	}

	// run task list
	var res string
	err = c.Run(ctxt, submit(`https://sit1.pressly.com/s3upgrade03/englishb-en-ca`,
		`ekrchan.sit1+1@gmail.com`, `******`, &res))
	if err != nil {
		log.Fatal(err)
	}

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("overview: %s", res)
}

func submit(urlstr, email string, pw string, res *string) cdp.Tasks {
	return cdp.Tasks{
		cdp.Navigate(urlstr),
		cdp.WaitVisible(`//input[@name="email"]`),
		cdp.SendKeys(`//input[@name="email"]`, email),
		cdp.SendKeys(`//input[@name="password"]`, pw),
		cdp.Click(`button[type="submit"]`),
		cdp.WaitNotPresent(`//input[@name="email"]`),
		cdp.Sleep(30 * time.Second),
		//cdp.Text(`//*[@id="js-pjax-container"]/div[2]/div/div[2]/ul/li/p`, res),
	}
}
