import no.elhub.devxp.build.configuration.pipeline.ElhubProject.Companion.elhubProject
import no.elhub.devxp.build.configuration.pipeline.constants.Group.DEVXP
import no.elhub.devxp.build.configuration.pipeline.jobs.makeVerify
<<<<<<< HEAD

<<<<<<< HEAD

elhubProject(DEVXP, "gh-dxp") {

    params {
        param("env.PATH", "\$PATH:/usr/local/go/bin:/usr/bin")
        param("env.GOROOT", "/usr/local/go")
    }

    pipeline(withReleaseVersion = false) {
=======
elhubProject(DEVXP, "devxp-jira-scripts") {

    params {
        param("env.PATH", "\$PATH:/opt/go/1.21.6/bin")
        param("env.GOROOT", "/opt/go/1.21.6")
    }

    pipeline {
>>>>>>> 0b2cd29 (Test local GO params)
=======

elhubProject(DEVXP, "devxp-jira-scripts") {

    params {
        param("env.PATH", "\$PATH:/opt/go/1.21.6/bin")
        param("env.GOROOT", "/opt/go/1.21.6")
    }

    pipeline {
>>>>>>> 0ba05c5 (Update TC settings)
        sequential {
            makeVerify {
                disableSonarScan = true
            }
        }
    }
}
