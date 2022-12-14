import React, { FC, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { SchemaForm, JSONSchema, initFormData, UISchema } from '@/components';
import type * as Type from '@/common/interface';
import { useToast } from '@/hooks';
import {
  getRequireAndReservedTag,
  postRequireAndReservedTag,
} from '@/services';
import { handleFormError } from '@/utils';

import '../index.scss';

const Legal: FC = () => {
  const { t } = useTranslation('translation', {
    keyPrefix: 'admin.write',
  });
  const Toast = useToast();
  // const updateSiteInfo = siteInfoStore((state) => state.update);

  const schema: JSONSchema = {
    title: t('page_title'),
    properties: {
      recommend_tags: {
        type: 'string',
        title: t('recommend_tags.label'),
        description: t('recommend_tags.text'),
      },
      required_tag: {
        type: 'boolean',
        title: t('required_tag.title'),
        label: t('required_tag.label'),
        description: t('required_tag.text'),
      },
      reserved_tags: {
        type: 'string',
        title: t('reserved_tags.label'),
        description: t('reserved_tags.text'),
      },
    },
  };
  const uiSchema: UISchema = {
    recommend_tags: {
      'ui:widget': 'textarea',
      'ui:options': {
        rows: 10,
      },
    },
    required_tag: {
      'ui:widget': 'switch',
    },
    reserved_tags: {
      'ui:widget': 'textarea',
      'ui:options': {
        rows: 10,
      },
    },
  };
  const [formData, setFormData] = useState(initFormData(schema));

  const onSubmit = (evt) => {
    evt.preventDefault();
    evt.stopPropagation();
    let recommend_tags = [];
    if (formData.recommend_tags.value?.trim()) {
      recommend_tags = formData.recommend_tags.value.trim().split('\n');
    }
    let reserved_tags = [];
    if (formData.reserved_tags.value?.trim()) {
      reserved_tags = formData.reserved_tags.value.trim().split('\n');
    }
    const reqParams: Type.AdminSettingsWrite = {
      recommend_tags,
      reserved_tags,
      required_tag: formData.required_tag.value,
    };
    postRequireAndReservedTag(reqParams)
      .then(() => {
        Toast.onShow({
          msg: t('update', { keyPrefix: 'toast' }),
          variant: 'success',
        });
      })
      .catch((err) => {
        if (err.isError) {
          const data = handleFormError(err, formData);
          setFormData({ ...data });
        }
      });
  };

  const initData = () => {
    getRequireAndReservedTag().then((res) => {
      if (Array.isArray(res.recommend_tags)) {
        formData.recommend_tags.value = res.recommend_tags.join('\n');
      }
      formData.required_tag.value = res.required_tag;
      if (Array.isArray(res.reserved_tags)) {
        formData.reserved_tags.value = res.reserved_tags.join('\n');
      }
      setFormData({ ...formData });
    });
  };

  useEffect(() => {
    initData();
  }, []);

  const handleOnChange = (data) => {
    setFormData(data);
  };

  return (
    <>
      <h3 className="mb-4">{t('page_title')}</h3>
      <SchemaForm
        schema={schema}
        formData={formData}
        onSubmit={onSubmit}
        uiSchema={uiSchema}
        onChange={handleOnChange}
      />
    </>
  );
};

export default Legal;
