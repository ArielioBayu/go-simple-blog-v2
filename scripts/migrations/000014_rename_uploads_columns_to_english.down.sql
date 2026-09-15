ALTER TABLE uploads
    CHANGE COLUMN file_name nama_file VARCHAR(255) NOT NULL,
    CHANGE COLUMN system_filename nama_file_sistem VARCHAR(255) NOT NULL,
    CHANGE COLUMN file_type type_file VARCHAR(50) NOT NULL;
