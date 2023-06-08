from faker import Faker
import pandas as pd
from tqdm import tqdm

def generate_and_save_data(num_records):
    faker = Faker()

    data = []
    total_size = 0  # Track the total data size

    # Generate the data and calculate the total size
    for i in tqdm(range(num_records), desc="Generating Data"):
        name = faker.name()
        token = faker.pystr(30, 30)
        data.append({
            "id": i + 1,
            "created_at": faker.date_time(),
            "updated_at": faker.date_time(),
            "deleted_at": faker.date_time(),
            "name": name,
            "token": token
        })
        total_size += len(name) + len(token)

    # Create a DataFrame from the list of dictionaries
    df = pd.DataFrame(data)

    # Print the DataFrame
    print(df.head())

    # Save DataFrame to CSV file
    with tqdm(total=total_size, unit='B', unit_scale=True, desc="Saving CSV") as pbar:
        df.to_csv('output1.csv', index=False, line_terminator='\n', chunksize=4096, encoding='utf-8')
        pbar.update(total_size)  # Update the progress bar to indicate completion

# Usage
generate_and_save_data(10000000)
