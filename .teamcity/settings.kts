import no.elhub.devxp.build.configuration.pipeline.ElhubProject.Companion.elhubProject
import no.elhub.devxp.build.configuration.pipeline.constants.Group.DEVXP
import no.elhub.devxp.build.configuration.pipeline.jobs.makeVerify


elhubProject(DEVXP, "gh-dxp-tc-build") {

    params {
        param("env.PATH", "%env.PATH%:/usr/local/go/bin:/usr/bin")
        param("env.GOROOT", "/usr/local/go")
    }

    pipeline(withReleaseVersion = false) {
        sequential {
            makeVerify {
                disableSonarScan = true
            }
        }
    }
}
