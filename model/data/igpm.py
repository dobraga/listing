from requests import get
from lxml import html
import pandas as pd


def read():
    complete_data = []

    r = get('http://www.yahii.com.br/igpm.html')
    r.raise_for_status()
    lines = html.fromstring(r.content).xpath(
        '//center/div[1]/center/table/tbody/tr')

    for line in lines[1::]:
        parsed_data = list(filter(lambda x: x != '', map(
            lambda x: x.strip(), line.xpath('.//text()'))))

        if not parsed_data[0].isdigit() or len(parsed_data) <= 1:
            continue

        for month, month_data in enumerate(parsed_data[1::], 1):
            if month > 12:
                continue

            complete_data.append({
                'year': parsed_data[0],
                'month': str(month),
                'data': month_data
            })

    data = pd.DataFrame(complete_data)

    data['sign'] = data['data'].str.contains(
        '(-)', regex=False).map({True: -1, False: 1})
    data['igpm'] = data['sign'] * data['data'].str.extract(
        r'([\d,]+)', expand=False).str.replace(',', '.', regex=False).astype(float)

    data['month'] = pd.to_datetime(
        data.year + '-' + data.month.str.zfill(2) + '-01')

    return data[['month', 'igpm']]


if __name__ == '__main__':
    print(read())
