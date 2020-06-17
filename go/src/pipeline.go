package src

import "errors"

type Pipeline struct {
	config  Config
	emailer Emailer
	log     Logger
}

func (p *Pipeline) run(project Project) {
	err := p.runTest(project)
	if err != nil {
		p.handleError(err)
		return
	}

	err = p.deploy(project)
	if err != nil {
		p.handleError(err)
		return
	}

	p.email("Deployment completed successfully")
}

func (p *Pipeline) handleError(err error) {
	p.log.error(err.Error())
	p.email(err.Error())
}

func (p *Pipeline) email(content string)  {
	if !p.config.sendEmailSummary() {
		p.log.info("Email disabled")
		return
	}

	p.log.info("Sending email")
	p.emailer.send(content)
}

func (p *Pipeline) deploy(project Project) error {
	if "success" == project.deploy() {
		p.log.info("Deployment successful")
		return nil
	}

	return errors.New("Deployment failed")
}

func (p *Pipeline) runTest(project Project) error {
	if !project.hasTests() {
		p.log.info("No tests")
		return nil
	}

	if "success" == project.runTests() {
		p.log.info("Tests passed")
		return nil
	}

	return errors.New("Tests failed")
}
