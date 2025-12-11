import { Alert, Button, Col, Container, Form, Row, Table } from "react-bootstrap";

import { useEffect, useState } from 'react';
import { Link } from 'react-router';
import { useDispatch, useSelector } from '../store';
import { BeatLoader } from "react-spinners";
import { ROUTE_LABELS, ROUTES } from "../routes";
import { Breadcrumbs } from "../components/breadcrumbs";
import { completeGenerationRequest, getGenerationRequests, rejectGenerationRequest, setFormedAtBegin, setFormedAtEnd, setStatus } from "../slices/generation-requests";

export const GenerationRequestsPage = () => {
    const dispatch = useDispatch();

    const [selectedId, setSelectedId] = useState<number | null>(null);

    const { loading, error, generationRequests, status, formedAtBegin, formedAtEnd } = useSelector(state => state.generationRequests);
    const { isEngineer } = useSelector(state => state.user);

    useEffect(() => {
        dispatch(getGenerationRequests());
    }, [status, formedAtBegin, formedAtEnd]);

    return (<Container>
        <Breadcrumbs breadcrumbs={[
            { label: ROUTE_LABELS.HOME, path: ROUTES.HOME },
            { label: `Заявки`, path: `` }
        ]} />
        <h1 style={{ marginBottom: 25 }}>Заявки</h1>
        {error && <Alert variant="danger">{error}</Alert>}
        <Row className="mb-4">
            <Col md={4}>
                <Form.Group>
                    <Form.Label>Статус</Form.Label>
                    <Form.Select
                        value={status ?? ""}
                        onChange={e => dispatch(setStatus(e.target.value || undefined as any))}
                    >
                        <option value="">Любой статус</option>
                        <option value="sent">В работе</option>
                        <option value="rejected">Отклонена</option>
                        <option value="completed">Завершена</option>
                    </Form.Select>
                </Form.Group>
            </Col>

            <Col md={4}>
                <Form.Group>
                    <Form.Label>Оформлено с</Form.Label>
                    <Form.Control
                        type="date"
                        value={formedAtBegin}
                        onChange={e => dispatch(setFormedAtBegin(e.target.value))}
                    />
                </Form.Group>
            </Col>

            <Col md={4}>
                <Form.Group>
                    <Form.Label>Оформлено до</Form.Label>
                    <Form.Control
                        type="date"
                        value={formedAtEnd}
                        onChange={e => dispatch(setFormedAtEnd(e.target.value))}
                    />
                </Form.Group>
            </Col>
        </Row>

        {loading
            ? <Container style={{ display: "flex", justifyContent: "center", paddingTop: "30%" }}>
                <BeatLoader color="#5BA1D4" />
            </Container>
            :
            <>
                <Table bordered hover>
                    <thead>
                        <tr>
                            <th>#</th>
                            <th>Статус</th>
                            <th>Дата создания</th>
                            <th>Дата оформления</th>
                            <th>Дата завершения</th>
                            {isEngineer && <>
                                <th>Создал</th>
                                <th>Обработал</th>
                            </>
                            }
                            <th>Турбин</th>
                        </tr>
                    </thead>
                    <tbody>
                        {generationRequests.map(generationRequest => (
                            <tr key={generationRequest.id} onClick={() => setSelectedId(generationRequest.id || -1)}>
                                <td>
                                    <Link to={`${ROUTES.GENERATION_REQUESTS}/${generationRequest.id}`}>
                                        {generationRequest.id}
                                    </Link>
                                </td>
                                <td>{{
                                    sent: "В работе",
                                    rejected: "Отклонена",
                                    completed: "Завершена"
                                }[generationRequest.status || ""]}</td>
                                <td>{generationRequest.created_at && (new Date(generationRequest.created_at ?? "")).toLocaleString()}</td>
                                <td>{generationRequest.formed_at && (new Date(generationRequest.formed_at ?? "")).toLocaleString()}</td>
                                <td>{generationRequest.closed_at && (new Date(generationRequest.closed_at ?? "")).toLocaleString()}</td>
                                {isEngineer &&
                                    <>
                                        <td>{generationRequest.created_by}</td>
                                        <td>{generationRequest.closed_by}</td>
                                    </>
                                }
                                <td>{generationRequest.turbine_generation_requests_count || "н/д"}</td>
                            </tr>
                        ))}
                    </tbody>
                </Table>
                {(generationRequests.find(generationRequest => generationRequest.id === selectedId)?.status === "sent") && (
                    <div
                        style={{
                            position: "fixed",
                            bottom: 0,
                            left: 0,
                            right: 0,
                            padding: "15px 20px",
                            backgroundColor: "white",
                            borderTop: "1px solid #ccc",
                            boxShadow: "0 -2px 6px rgba(0,0,0,0.1)",
                            zIndex: 999,
                        }}
                    >
                        <Container
                            style={{
                                display: "flex",
                                justifyContent: "space-between",
                                alignItems: "center",
                            }}
                        >
                            <strong>
                                Заявка №{selectedId}
                            </strong>

                            <div className="d-flex gap-2">
                                <Button variant="success" onClick={() => { dispatch(completeGenerationRequest(selectedId as number)); dispatch(getGenerationRequests()); }}>Завершить</Button>
                                <Button variant="danger" onClick={() => { dispatch(rejectGenerationRequest(selectedId as number)); dispatch(getGenerationRequests()); }}>Отклонить</Button>
                            </div>
                        </Container>
                    </div>
                )}
            </>
        }</Container>)
};