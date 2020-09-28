import React, {Component} from 'react';
import GraphiQL from 'graphiql';
import fetch from 'isomorphic-fetch';
import './App.css'

class App extends Component {
  constructor(props){
    super(props)
    this.state = {
      authenticated: false
    }
  }

  componentDidMount() {
    return fetch("http://api.firstcontributions.com/v1/session", {
      credentials: 'include'
    }).then(res=>{
      if (res.status < 400) {
        this.setState({authenticated: true})
      }
    })
  }

  graphQLFetcher(graphQLParams) {
    return fetch("http://api.firstcontributions.com/v1/graphql", {
      credentials: 'include',
      method: 'post',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(graphQLParams),
    }).then(response => response.json());
  }

  render() {
    const {authenticated} = this.state
    if (authenticated) {
      return (
        <div className="app">
          <GraphiQL fetcher={this.graphQLFetcher} />
        </div>
      )
    }
    return (
      <div>
        <a href={encodeURI("http://api.firstcontributions.com/v1/auth/redirect")}>
        <button>Login With Github</button>
        </a>
      </div>
    )
  }
}

export default App;
