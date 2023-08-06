from diagrams import Cluster, Diagram, Edge
from diagrams.onprem.client import Users
from diagrams.onprem.monitoring import Prometheus
from diagrams.programming.language import Go


graph_attr = {
    'bgcolor': 'white'
}

with Diagram('Prometheus setup in Docker Compose', show=False, graph_attr=graph_attr):
    user = Users('Users')

    with Cluster('Docker Compose Orchestration'):
        prometheus = Prometheus('Prometheus')
        app = Go('Golang App')

        user >> Edge(label='http://localhost:8080/v1/save') >> app
        user >> Edge(label='http://localhost:8080/v1/spend') >> app

        app >> prometheus
