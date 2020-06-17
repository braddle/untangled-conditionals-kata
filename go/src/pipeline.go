package src

type Pipeline struct {
	config  Config
	emailer Emailer
	log     Logger
}

func (p *Pipeline) run(project Project) {
	testsPassed := p.runTest(project)

	deploySuccessful := p.deploy(project, testsPassed)

	p.emailResults(testsPassed, deploySuccessful)
}

func (p *Pipeline) emailResults(testsPassed bool, deploySuccessful bool) {
	if !testsPassed {
		p.email("Tests failed")
		return
	}

	if deploySuccessful {
		p.email("Deployment completed successfully")
		return
	}

	p.email("Deployment failed")
}

func (p *Pipeline) email(content string)  {
	if !p.config.sendEmailSummary() {
		p.log.info("Email disabled")
		return
	}

	p.log.info("Sending email")
	p.emailer.send(content)
}

func (p *Pipeline) deploy(project Project, testsPassed bool) bool {
	if !testsPassed {
		return false
	}

	if "success" == project.deploy() {
		p.log.info("Deployment successful")
		return true
	}

	p.log.error("Deployment failed")
	return false
}

func (p *Pipeline) runTest(project Project) bool {
	if !project.hasTests() {
		p.log.info("No tests")
		return true
	}

	if "success" == project.runTests() {
		p.log.info("Tests passed")
		return true
	}

	p.log.error("Tests failed")
	return false
}
