import no.elhub.devxp.build.configuration.pipeline.ElhubProject.Companion.elhubProject
import no.elhub.devxp.build.configuration.pipeline.constants.Group.DEVXP
import no.elhub.devxp.build.configuration.pipeline.jobs.makeVerify

elhubProject(DEVXP, "devxp-jira-scripts") {

    params {
        param("env.PATH", "/opt/go/1.21.6/bin:%teamcity.tool.maven.DEFAULT%/bin:%KOTLIN_PATH%:%env.PATH%:/home/teamcity/.nvm/versions/node/v20.10.0/bin")
        param("env.GOROOT", "/opt/go/1.21.6")
    }

    pipeline {
        sequential {
            makeVerify {
                disableSonarScan = true
            }
        }
    }
}
