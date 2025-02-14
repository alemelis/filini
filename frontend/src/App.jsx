import { useState } from "react";
import reactLogo from "./assets/react.svg";
import viteLogo from "/vite.svg";
import "./App.css";
import Search from "./Search";

function App() {
  const [count, setCount] = useState(0);

  return (
    <>
      <div className="App">
        <h1 className="text-center text-3xl mt-6">Filini Subtitle Search</h1>
        <Search /> {/* Add the Search component */}
      </div>
    </>
  );
}

export default App;
