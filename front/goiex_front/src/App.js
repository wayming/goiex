import React, {useState, useEffect} from 'react';
import {Button, Form, Container, Modal } from 'react-bootstrap'
import axios from "axios";

import './App.css';
import Symbol from './symbol.js'

function App() {

  const [symbols, setSymbols] = useState([])

  return (
    <div>
      <Container>
        <Button onClick={() => loadSymbols()}>Load Symbols</Button>
        <Button onClick={() => getSymbols()}>Get Symbols</Button>
      </Container>

      <Container>
          {symbols != null && symbols.map((symbol, i) => (
              <Symbol symbolData={symbol}/>
          ))}
      </Container>
    </div>

  );

  function loadSymbols(){
    var url = "/load"
    axios.get(url, { responseType: 'json' }).then(response => {
        console.log(response.status)
        if(response.status == 200){
          setSymbols(response.data)
        }
    })
  }
    
  function getSymbols(){
    var url = "/symbols"
    axios.get(url, { responseType: 'json' }).then(response => {
        console.log(response.status)
        if(response.status == 200){
          setSymbols(response.data)
        }
    })
  }
}

export default App;
