import no.elhub.devxp.build.configuration.pipeline.ElhubProject.Companion.elhubProject
import no.elhub.devxp.build.configuration.pipeline.constants.Group.DEVXP
import no.elhub.devxp.build.configuration.pipeline.jobs.makeVerify


elhubProject(DEVXP, "gh-dxp") {

    params {
        param("env.PATH", "\$PATH:/opt/go/1.21.6/bin")
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
