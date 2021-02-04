library "alauda-cicd"
def language = "golang"
AlaudaPipeline {
    config = [
        agent: 'golang-1.13',
        folder: '.',
//         chart: [
//             [
//                 project: "timatrix",
//                 pipeline: "ti-installer-chart-update",
//                 component: "timatrixController",
//             ],
//         ],
        scm: [
            credentials: 'alaudabot-bitbucket'
        ],
        docker: [
            repository: "asm/operator-monitor",
            context: ".",
            dockerfile: "Dockerfile",
            armBuild: false,
        ],
        sonar: [
            binding: "sonarqube",
// 			enabled: false
        ],
    ]
    env = [
        GOPROXY: "https://athens.alauda.cn",
//         CGO_ENABLED: "0",
//         GOOS: "linux",
    ]
    yaml = "alauda.yaml"
}
