import {
  createBrowserRouter,
  RouterProvider,
  Outlet,
  Route,
  Link,
} from "react-router-dom";
import Dashboard from "./pages/Dashboard";
import Login from "./pages/Login";
import { useEffect } from "react";
import TotpFallbackPage from "./pages/TotpFallbackPage/TotpFallbackPage";

const Root = ({ children }) => {
  return <Outlet>{children}</Outlet>;
};

const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
    children: [
      { path: "", element: <Login /> },
      { path: "dashboard", element: <Dashboard /> },
      { path: "totp/:username", element: <TotpFallbackPage /> },
    ],
  },
]);

function App() {
  useEffect(() => {
      // check whether current browser supports WebAuthn
      if (!window.PublicKeyCredential) {
          console.log("doesn't support supports webauthn");
          alert('Error: this browser does not support WebAuthn');
          return;
      }
  }, [])
  
  return (
    <div>
      <RouterProvider router={router} />
    </div>
  )
}

export default App
