import { useState } from "react";
import axios from "axios";

function Search() {
  const [query, setQuery] = useState(""); // store the search query
  const [results, setResults] = useState([]); // store the results

  const handleSearch = async () => {
    try {
      // Make a GET request to the backend's search endpoint
      const response = await axios.get(
        `http://localhost:8080/subtitles/search?q=${query}`,
      );
      setResults(response.data); // update the results state with data from the backend
    } catch (error) {
      console.error("Error fetching search results:", error);
    }
  };

  const handleKeyDown = (e) => {
    if (e.key === "Enter") {
      handleSearch(); // Trigger search when Enter is pressed
    }
  };

  return (
    <div className="p-4">
      {/* Search Input */}
      <input
        type="text"
        placeholder="Search quote..."
        className="p-2 border border-gray-300"
        value={query}
        onChange={(e) => setQuery(e.target.value)} // update query when user types
        onKeyDown={handleKeyDown}
      />

      {/* Search Button */}
      <button
        onClick={handleSearch}
        className="bg-blue-500 text-white p-2 ml-4"
      >
        Search
      </button>

      {/* Display Search Results */}
      <ul className="mt-4">
        {results.map((item, index) => (
          <li key={index} className="border p-2 mt-2">
            <div>{item.quote}</div>
            <div>
              <img
                src={`http://localhost:8080/${item.file_path}`}
                alt="Subtitle GIF"
              />
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default Search;
