from diagrams import Cluster, Diagram, Edge
from diagrams.generic.device import Mobile
from diagrams.onprem.client import Users
from diagrams.onprem.monitoring import Prometheus
from diagrams.onprem.compute import Server

graph_attr = {
    'bgcolor': 'white'
}

with Diagram('Prometheus setup in Docker Compose', show=False, graph_attr=graph_attr):
    user = Users('Users')
    node = Server('Server Node')
    device = Mobile('On-call Engineer')

    with Cluster('Docker Compose Orchestration'):
        prometheus = Prometheus('Prometheus')
        alertmanager = Prometheus('Alertmanager')
        node_exporter = Prometheus('Node exporter')

        user >> Edge(label='Port: 9090') >> prometheus

        prometheus >> node_exporter >> node
        prometheus >> alertmanager >> device
