import React, { useEffect, useState } from 'react';
import logo from './logo.svg';
import apiUrl from './helpers/apiUrl';
import './App.css';
import axios from 'axios';

function App() {
  const [helloResponse, setHelloResponse] = useState({});

  useEffect(() => {
    getHello()
  })

  function getHello() {
    axios
      .get(apiUrl("/api/hello"), {
        headers: {
          Accept: "application/json",
        },
      })
      .then((response) => {
        setHelloResponse(response.data)
      });
  }

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <span>
          <strong>/api/hello</strong> returned:
        </span>
        <pre>{JSON.stringify(helloResponse)}</pre>
      </header>
    </div>
  );
}

export default App;
