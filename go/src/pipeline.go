package src

type Pipeline struct {
	config  Config
	emailer Emailer
	log     Logger
}

func (p *Pipeline) run(project Project) {
	var deploySuccessful = false

	testsPassed := p.runTest(project)

	if testsPassed {
		deploySuccessful = p.deploy(project)
	}

	if p.config.sendEmailSummary() {
		p.log.info("Sending email")
		if testsPassed {
			if deploySuccessful {
				p.emailer.send("Deployment completed successfully")
			} else {
				p.emailer.send("Deployment failed")
			}
		} else {
			p.emailer.send("Tests failed")
		}
	} else {
		p.log.info("Email disabled")
	}
}

func (p *Pipeline) deploy(project Project) bool {
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
	} else {
		p.log.error("Tests failed")
		return false
	}
}
