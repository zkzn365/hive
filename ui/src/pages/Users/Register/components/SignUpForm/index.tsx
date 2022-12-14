import React, { FormEvent, MouseEvent, useState } from 'react';
import { Form, Button, Col } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import { Trans, useTranslation } from 'react-i18next';

import type { FormDataType } from '@/common/interface';
import { register, useLegalTos, useLegalPrivacy } from '@/services';
import userStore from '@/stores/userInfo';
import { handleFormError } from '@/utils';

interface Props {
  callback: () => void;
}

const Index: React.FC<Props> = ({ callback }) => {
  const { t } = useTranslation('translation', { keyPrefix: 'login' });
  const [formData, setFormData] = useState<FormDataType>({
    name: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
    e_mail: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
    pass: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
    captcha_code: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
  });
  const updateUser = userStore((state) => state.update);

  const handleChange = (params: FormDataType) => {
    setFormData({ ...formData, ...params });
  };

  const checkValidated = (): boolean => {
    let bol = true;
    const { name, e_mail, pass } = formData;
    if (!name.value) {
      bol = false;
      formData.name = {
        value: '',
        isInvalid: true,
        errorMsg: t('name.msg.empty'),
      };
    } else if ([...name.value].length > 30) {
      bol = false;
      formData.name = {
        value: name.value,
        isInvalid: true,
        errorMsg: t('name.msg.range'),
      };
    }

    if (!e_mail.value) {
      bol = false;
      formData.e_mail = {
        value: '',
        isInvalid: true,
        errorMsg: t('email.msg.empty'),
      };
    }

    if (!pass.value) {
      bol = false;
      formData.pass = {
        value: '',
        isInvalid: true,
        errorMsg: t('password.msg.empty'),
      };
    }
    setFormData({
      ...formData,
    });
    return bol;
  };
  const { data: tos } = useLegalTos();
  const { data: privacy } = useLegalPrivacy();
  const argumentClick = (evt: MouseEvent, type: 'tos' | 'privacy') => {
    evt.stopPropagation();
    const contentText =
      type === 'tos'
        ? tos?.terms_of_service_original_text
        : privacy?.privacy_policy_original_text;
    let matchUrl: URL | undefined;
    try {
      if (contentText) {
        matchUrl = new URL(contentText);
      }
      // eslint-disable-next-line no-empty
    } catch (ex) {}
    if (matchUrl) {
      evt.preventDefault();
      window.open(matchUrl.toString());
    }
  };

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    event.stopPropagation();
    if (!checkValidated()) {
      return;
    }
    register({
      name: formData.name.value,
      e_mail: formData.e_mail.value,
      pass: formData.pass.value,
    })
      .then((res) => {
        updateUser(res);
        callback();
      })
      .catch((err) => {
        if (err.isError) {
          const data = handleFormError(err, formData);
          setFormData({ ...data });
        }
      });
  };

  return (
    <Col className="mx-auto" md={3}>
      <Form noValidate onSubmit={handleSubmit} autoComplete="off">
        <Form.Group controlId="name" className="mb-3">
          <Form.Label>{t('name.label')}</Form.Label>
          <Form.Control
            autoComplete="off"
            required
            type="text"
            isInvalid={formData.name.isInvalid}
            value={formData.name.value}
            onChange={(e) =>
              handleChange({
                name: {
                  value: e.target.value,
                  isInvalid: false,
                  errorMsg: '',
                },
              })
            }
          />
          <Form.Control.Feedback type="invalid">
            {formData.name.errorMsg}
          </Form.Control.Feedback>
        </Form.Group>
        <Form.Group controlId="email" className="mb-3">
          <Form.Label>{t('email.label')}</Form.Label>
          <Form.Control
            autoComplete="off"
            required
            type="e_mail"
            isInvalid={formData.e_mail.isInvalid}
            value={formData.e_mail.value}
            onChange={(e) =>
              handleChange({
                e_mail: {
                  value: e.target.value,
                  isInvalid: false,
                  errorMsg: '',
                },
              })
            }
          />
          <Form.Control.Feedback type="invalid">
            {formData.e_mail.errorMsg}
          </Form.Control.Feedback>
        </Form.Group>

        <Form.Group controlId="password" className="mb-3">
          <Form.Label>{t('password.label')}</Form.Label>
          <Form.Control
            autoComplete="off"
            required
            type="password"
            maxLength={32}
            isInvalid={formData.pass.isInvalid}
            value={formData.pass.value}
            onChange={(e) =>
              handleChange({
                pass: {
                  value: e.target.value,
                  isInvalid: false,
                  errorMsg: '',
                },
              })
            }
          />
          <Form.Control.Feedback type="invalid">
            {formData.pass.errorMsg}
          </Form.Control.Feedback>
        </Form.Group>

        <div className="d-grid">
          <Button variant="primary" type="submit">
            {t('signup', { keyPrefix: 'btns' })}
          </Button>
        </div>
      </Form>
      <div className="text-center fs-14 mt-3">
        <Trans i18nKey="login.agreements" ns="translation">
          By registering, you agree to the
          <Link
            to="/privacy"
            onClick={(evt) => {
              argumentClick(evt, 'privacy');
            }}
            target="_blank">
            privacy policy
          </Link>
          and
          <Link
            to="/tos"
            onClick={(evt) => {
              argumentClick(evt, 'tos');
            }}
            target="_blank">
            terms of service
          </Link>
          .
        </Trans>
      </div>
      <div className="text-center mt-5">
        <Trans i18nKey="login.info_login" ns="translation">
          Already have an account? <Link to="/users/login">Log in</Link>
        </Trans>
      </div>
    </Col>
  );
};

export default React.memo(Index);
