/* eslint-disable prettier/prettier */
import { FC, useState, useEffect } from 'react';
import { Container, Row, Col, Card, Alert } from 'react-bootstrap';
import { useTranslation, Trans } from 'react-i18next';

import type { FormDataType } from '@/common/interface';
import { PageTitle } from '@/components';
import {
  dbCheck,
  installInit,
  installBaseInfo,
  checkConfigFileExists,
} from '@/services';
import { Storage, handleFormError } from '@/utils';
import { CURRENT_LANG_STORAGE_KEY } from '@/common/constants';

import {
  FirstStep,
  SecondStep,
  ThirdStep,
  FourthStep,
  Fifth,
} from './components';

const Index: FC = () => {
  const { t } = useTranslation('translation', { keyPrefix: 'install' });
  const [step, setStep] = useState(1);
  const [loading, setLoading] = useState(true);
  const [errorData, setErrorData] = useState<{ [propName: string]: any }>({
    msg: '',
  });
  const [checkData, setCheckData] = useState({
    db_table_exist: false,
    db_connection_success: false,
  });

  const [formData, setFormData] = useState<FormDataType>({
    lang: {
      value: 'en_US',
      isInvalid: false,
      errorMsg: '',
    },
    db_type: {
      value: 'mysql',
      isInvalid: false,
      errorMsg: '',
    },
    db_username: {
      value: 'root',
      isInvalid: false,
      errorMsg: '',
    },
    db_password: {
      value: 'root',
      isInvalid: false,
      errorMsg: '',
    },
    db_host: {
      value: 'db:3306',
      isInvalid: false,
      errorMsg: '',
    },
    db_name: {
      value: 'answer',
      isInvalid: false,
      errorMsg: '',
    },
    db_file: {
      value: '/data/answer.db',
      isInvalid: false,
      errorMsg: '',
    },
    site_name: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
    site_url: {
      value: window.location.origin,
      isInvalid: false,
      errorMsg: '',
    },
    contact_email: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
    name: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
    password: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
    email: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
  });

  const handleChange = (params: FormDataType) => {
    // console.log(params);
    setErrorData({
      msg: '',
    });
    setFormData({ ...formData, ...params });
  };

  const handleErr = (data) => {
    window.scrollTo(0, 0);
    setErrorData(data);
  };

  const handleNext = async () => {
    setErrorData({
      msg: '',
    });
    setStep((pre) => pre + 1);
  };

  const checkInstall = () => {
    const params = {
      lang: formData.lang.value,
      db_type: formData.db_type.value,
      db_username: formData.db_username.value,
      db_password: formData.db_password.value,
      db_host: formData.db_host.value,
      db_name: formData.db_name.value,
      db_file: formData.db_file.value,
    };
    installInit(params)
      .then(() => {
        handleNext();
      })
      .catch((err) => {
        handleErr(err);
      });
  };

  const submitDatabaseForm = () => {
    const params = {
      lang: formData.lang.value,
      db_type: formData.db_type.value,
      db_username: formData.db_username.value,
      db_password: formData.db_password.value,
      db_host: formData.db_host.value,
      db_name: formData.db_name.value,
      db_file: formData.db_file.value,
    };
    dbCheck(params)
      .then(() => {
        checkInstall();
      })
      .catch((err) => {
        handleErr(err);
      });
  };

  const submitSiteConfig = () => {
    const params = {
      lang: formData.lang.value,
      site_name: formData.site_name.value,
      site_url: formData.site_url.value,
      contact_email: formData.contact_email.value,
      name: formData.name.value,
      password: formData.password.value,
      email: formData.email.value,
    };
    installBaseInfo(params)
      .then(() => {
        handleNext();
      })
      .catch((err) => {
        if (err.isError) {
          const data = handleFormError(err, formData);
          setFormData({ ...data });
        } else {
          handleErr(err);
        }
      });
  };

  const handleStep = () => {
    if (step === 1) {
      Storage.set(CURRENT_LANG_STORAGE_KEY, formData.lang.value);
      handleNext();
    }
    if (step === 2) {
      submitDatabaseForm();
    }
    if (step === 3) {
      if (errorData.msg) {
        checkInstall();
      } else {
        handleNext();
      }
    }
    if (step === 4) {
      submitSiteConfig();
    }
    if (step > 4) {
      handleNext();
    }
  };

  const handleInstallNow = (e) => {
    e.preventDefault();
    if (checkData.db_table_exist) {
      setStep(8);
    } else {
      setStep(4);
    }
  };

  const configYmlCheck = () => {
    checkConfigFileExists()
      .then((res) => {
        setCheckData({
          db_table_exist: res.db_table_exist,
          db_connection_success: res.db_connection_success,
        });
        if (res && res.config_file_exist) {
          if (res.db_connection_success) {
            setStep(6)
          } else {
            setStep(7);
          }
        }
      })
      .finally(() => {
        setLoading(false);
      });
  };

  useEffect(() => {
    configYmlCheck();
  }, []);

  if (loading) {
    return <div />;
  }

  return (
    <div className="bg-f5 py-5 flex-grow-1">
      <PageTitle title={t('install', { keyPrefix: 'page_title' })} />
      <Container className='py-3'>
        <Row className="justify-content-center">
          <Col lg={6}>
            <h2 className="mb-4 text-center">{t('title')}</h2>
            <Card>
              <Card.Body>
                {errorData?.msg && (
                  <Alert variant="danger">{errorData?.msg}</Alert>
                )}

                <FirstStep
                  visible={step === 1}
                  data={formData.lang}
                  changeCallback={handleChange}
                  nextCallback={handleStep}
                />

                <SecondStep
                  visible={step === 2}
                  data={formData}
                  changeCallback={handleChange}
                  nextCallback={handleStep}
                />

                <ThirdStep
                  visible={step === 3}
                  nextCallback={handleStep}
                  errorMsg={errorData}
                />

                <FourthStep
                  visible={step === 4}
                  data={formData}
                  changeCallback={handleChange}
                  nextCallback={handleStep}
                />

                <Fifth visible={step === 5} siteUrl={formData.site_url.value} />
                {step === 6 && (
                  <div>
                    <h5>{t('warn_title')}</h5>
                    <p>
                      <Trans i18nKey="install.warn_description" components={{ 1: <code />}} />
                      {' '}
                      <Trans i18nKey="install.install_now">
                        You may try <a href="###" onClick={e => handleInstallNow(e)}>installing now</a>.
                      </Trans>
                    </p>
                  </div>
                )}

                {step === 7 && (
                  <div>
                    <h5>{t('db_failed')}</h5>
                    <p>
                      <Trans i18nKey="install.db_failed_description" components={{ 1: <code />}} />
                    </p>
                  </div>
                )}

                {step === 8  && (
                  <div>
                    <h5>{t('installed')}</h5>
                    <p>{t('installed_description')}</p>
                  </div>
                )}
              </Card.Body>
            </Card>
          </Col>
        </Row>
      </Container>
    </div>
  );
};

export default Index;
