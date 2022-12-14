import { FC, FormEvent } from 'react';
import { Form, Button } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';

import Progress from '../Progress';
import type { FormDataType } from '@/common/interface';

interface Props {
  data: FormDataType;
  changeCallback: (value: FormDataType) => void;
  nextCallback: () => void;
  visible: boolean;
}

const sqlData = [
  {
    value: 'mysql',
    label: 'MariaDB/MySQL',
  },
  {
    value: 'sqlite3',
    label: 'SQLite',
  },
  {
    value: 'postgres',
    label: 'PostgreSQL',
  },
];

const Index: FC<Props> = ({ visible, data, changeCallback, nextCallback }) => {
  const { t } = useTranslation('translation', { keyPrefix: 'install' });

  const checkValidated = (): boolean => {
    let bol = true;
    const { db_type, db_username, db_password, db_host, db_name, db_file } =
      data;

    if (db_type.value !== 'sqlite3') {
      if (!db_username.value) {
        bol = false;
        data.db_username = {
          value: '',
          isInvalid: true,
          errorMsg: t('db_username.msg'),
        };
      }

      if (!db_password.value) {
        bol = false;
        data.db_password = {
          value: '',
          isInvalid: true,
          errorMsg: t('db_password.msg'),
        };
      }

      if (!db_host.value) {
        bol = false;
        data.db_host = {
          value: '',
          isInvalid: true,
          errorMsg: t('db_host.msg'),
        };
      }

      if (!db_name.value) {
        bol = false;
        data.db_name = {
          value: '',
          isInvalid: true,
          errorMsg: t('db_name.msg'),
        };
      }
    } else if (!db_file.value) {
      bol = false;
      data.db_file = {
        value: '',
        isInvalid: true,
        errorMsg: t('db_file.msg'),
      };
    }
    changeCallback({
      ...data,
    });
    return bol;
  };

  const handleSubmit = (event: FormEvent) => {
    event.preventDefault();
    event.stopPropagation();
    if (!checkValidated()) {
      return;
    }
    nextCallback();
  };

  if (!visible) return null;
  return (
    <Form noValidate onSubmit={handleSubmit}>
      <Form.Group controlId="database_engine" className="mb-3">
        <Form.Label>{t('db_type.label')}</Form.Label>
        <Form.Select
          value={data.db_type.value}
          isInvalid={data.db_type.isInvalid}
          onChange={(e) => {
            changeCallback({
              db_type: {
                value: e.target.value,
                isInvalid: false,
                errorMsg: '',
              },
            });
          }}>
          {sqlData.map((item) => {
            return (
              <option key={item.value} value={item.value}>
                {item.label}
              </option>
            );
          })}
        </Form.Select>
      </Form.Group>
      {data.db_type.value !== 'sqlite3' ? (
        <>
          <Form.Group controlId="username" className="mb-3">
            <Form.Label>{t('db_username.label')}</Form.Label>
            <Form.Control
              required
              placeholder={t('db_username.placeholder')}
              value={data.db_username.value}
              isInvalid={data.db_username.isInvalid}
              onChange={(e) => {
                changeCallback({
                  db_username: {
                    value: e.target.value,
                    isInvalid: false,
                    errorMsg: '',
                  },
                });
              }}
            />
            <Form.Control.Feedback type="invalid">
              {data.db_username.errorMsg}
            </Form.Control.Feedback>
          </Form.Group>

          <Form.Group controlId="db_password" className="mb-3">
            <Form.Label>{t('db_password.label')}</Form.Label>
            <Form.Control
              required
              value={data.db_password.value}
              isInvalid={data.db_password.isInvalid}
              onChange={(e) => {
                changeCallback({
                  db_password: {
                    value: e.target.value,
                    isInvalid: false,
                    errorMsg: '',
                  },
                });
              }}
            />

            <Form.Control.Feedback type="invalid">
              {data.db_password.errorMsg}
            </Form.Control.Feedback>
          </Form.Group>

          <Form.Group controlId="db_host" className="mb-3">
            <Form.Label>{t('db_host.label')}</Form.Label>
            <Form.Control
              required
              placeholder={t('db_host.placeholder')}
              value={data.db_host.value}
              isInvalid={data.db_host.isInvalid}
              onChange={(e) => {
                changeCallback({
                  db_host: {
                    value: e.target.value,
                    isInvalid: false,
                    errorMsg: '',
                  },
                });
              }}
            />
            <Form.Control.Feedback type="invalid">
              {data.db_host.errorMsg}
            </Form.Control.Feedback>
          </Form.Group>

          <Form.Group controlId="name" className="mb-3">
            <Form.Label>{t('db_name.label')}</Form.Label>
            <Form.Control
              required
              placeholder={t('db_name.placeholder')}
              value={data.db_name.value}
              isInvalid={data.db_name.isInvalid}
              onChange={(e) => {
                changeCallback({
                  db_name: {
                    value: e.target.value,
                    isInvalid: false,
                    errorMsg: '',
                  },
                });
              }}
            />
            <Form.Control.Feedback type="invalid">
              {data.db_name.errorMsg}
            </Form.Control.Feedback>
          </Form.Group>
        </>
      ) : (
        <Form.Group controlId="file" className="mb-3">
          <Form.Label>{t('db_file.label')}</Form.Label>
          <Form.Control
            required
            placeholder={t('db_file.placeholder')}
            value={data.db_file.value}
            isInvalid={data.db_file.isInvalid}
            onChange={(e) => {
              changeCallback({
                db_file: {
                  value: e.target.value,
                  isInvalid: false,
                  errorMsg: '',
                },
              });
            }}
          />
          <Form.Control.Feedback type="invalid">
            {data.db_file.errorMsg}
          </Form.Control.Feedback>
        </Form.Group>
      )}

      <div className="d-flex align-items-center justify-content-between">
        <Progress step={2} />
        <Button type="submit">{t('next')}</Button>
      </div>
    </Form>
  );
};

export default Index;
