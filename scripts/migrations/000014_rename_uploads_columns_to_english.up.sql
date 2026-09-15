ALTER TABLE uploads
    CHANGE COLUMN nama_file file_name VARCHAR(255) NOT NULL,
    CHANGE COLUMN nama_file_sistem system_filename VARCHAR(255) NOT NULL,
    CHANGE COLUMN type_file file_type VARCHAR(50) NOT NULL;
