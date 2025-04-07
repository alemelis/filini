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
      console.log("Search results:", response.data); // Log results for debugging
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
    <div className="container mx-auto p-4">
      {/* Search Input Container */}
      <div className="flex flex-wrap mb-6">
        <input
          type="text"
          placeholder="Search quote..."
          className="p-2 border border-gray-300 rounded flex-grow max-w-full"
          value={query}
          onChange={(e) => setQuery(e.target.value)} // update query when user types
          onKeyDown={handleKeyDown}
        />

        <button
          onClick={handleSearch}
          className="bg-blue-500 text-white p-2 ml-2 rounded"
        >
          Search
        </button>
      </div>

      {/* Results Count */}
      {results.length > 0 && (
        <div className="mb-4 text-gray-600">
          Found {results.length} result{results.length !== 1 ? "s" : ""}
        </div>
      )}

      {/* Display Search Results */}
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 w-full">
        {results.length > 0 ? (
          results.map((item, index) => (
            <div
              key={index}
              className="border p-3 rounded-lg shadow-sm bg-white flex flex-col"
            >
              <div className="font-medium mb-2 text-gray-800">
                {item.text || item.quote}
              </div>
              <div className="flex-grow">
                <video
                  src={`http://localhost:8080/${item.file_path}`}
                  autoPlay
                  loop
                  muted
                  playsInline
                  className="w-full h-auto rounded"
                />
                <div className="mt-2 flex justify-between items-center space-x-2">
                  <button
                    onClick={() => {
                      const videoUrl = `http://localhost:8080/${item.file_path}`;
                      navigator.clipboard.writeText(videoUrl);
                      alert("Video URL copied to clipboard!");
                    }}
                    className="bg-blue-500 hover:bg-blue-600 text-white text-sm py-1 px-2 rounded flex-1"
                    type="button"
                  >
                    Copy Link
                  </button>
                  <a
                    href={`http://localhost:8080/${item.file_path}`}
                    download
                    className="bg-green-500 hover:bg-green-600 text-white text-sm py-1 px-2 rounded flex-1 text-center"
                    role="button"
                  >
                    Download
                  </a>
                </div>
              </div>
            </div>
          ))
        ) : (
          <div className="col-span-full text-center py-8 text-gray-500">
            {query
              ? "No results found. Try a different search term."
              : "Enter a search term to find quotes."}
          </div>
        )}
      </div>
    </div>
  );
}

export default Search;
