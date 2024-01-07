CREATE TABLE IF NOT EXISTS collection (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  color VARCHAR(7) NOT NULL
);

CREATE TABLE IF NOT EXISTS todo (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  collection_id INT NOT NULL,
  text TEXT NOT NULL,
  done BOOLEAN DEFAULT FALSE,
  CONSTRAINT fk_collection
    FOREIGN KEY(collection_id) 
	    REFERENCES collection(id)
	    ON DELETE CASCADE
);