import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './pages/App';
import reportWebVitals from './reportWebVitals';
import {BrowserRouter, Route, Routes} from "react-router-dom";
import {ChakraProvider, ColorModeScript} from "@chakra-ui/react";
import Login from "./pages/Login";
import theme from "./theme"
import AuthTest from "./pages/AuthTest";
import SignUpPage from "./pages/SignUp";

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  <React.StrictMode>
    <ChakraProvider>
      <ColorModeScript initialColorMode={theme.config.initialColorMode}/>
      <Router/>
    </ChakraProvider>
  </React.StrictMode>
);

function Router() {
  return (
    <BrowserRouter basename={'/app'}>
      <Routes>
        <Route path="/" element={<App/>}/>
        <Route path={"/login"} element={<Login/>}/>
        <Route path={"/signup"} element={<SignUpPage/>}/>
        <Route path={"/authenticated"} element={<AuthTest/>}/>
      </Routes>
    </BrowserRouter>
  )
}

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
